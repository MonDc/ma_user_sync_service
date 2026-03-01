package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gorilla/mux"
    "github.com/jmoiron/sqlx"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    
    "github.com/yourusername/user-sync-service/internal/application/commands"
    "github.com/yourusername/user-sync-service/internal/application/queries"
    "github.com/yourusername/user-sync-service/internal/domain/user"
    "github.com/yourusername/user-sync-service/internal/infrastructure/api/handlers"
    "github.com/yourusername/user-sync-service/internal/infrastructure/config"
    "github.com/yourusername/user-sync-service/internal/infrastructure/observability/logger"
    "github.com/yourusername/user-sync-service/internal/infrastructure/observability/metrics"
    "github.com/yourusername/user-sync-service/internal/infrastructure/observability/tracing"
    "github.com/yourusername/user-sync-service/internal/infrastructure/persistence/mysql"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }

    // Initialize logger
    log, err := logger.NewLogger(cfg.Observability.LogLevel)
    if err != nil {
        log.Fatal("Failed to initialize logger:", err)
    }
    defer log.Sync()

    // Initialize tracing if enabled
    if cfg.Observability.Tracing.Enabled {
        cleanup, err := tracing.InitTracer(
            cfg.Observability.Tracing.ServiceName,
            cfg.Observability.Tracing.Environment,
            cfg.Observability.Tracing.ExporterURL,
        )
        if err != nil {
            log.Fatal("Failed to initialize tracer:", err)
        }
        defer cleanup()
    }

    // Initialize metrics
    metrics := metrics.NewMetrics("user_sync")

    // Initialize database connections
    mainDB, err := connectDatabase(cfg.Database.Main.Host, cfg.Database.Main.Port,
        cfg.Database.Main.User, cfg.Database.Main.Password, cfg.Database.Main.DBName,
        cfg.Database.Main.MaxConns, cfg.Database.Main.MinConns)
    if err != nil {
        log.Fatal("Failed to connect to main database:", err)
    }
    defer mainDB.Close()

    localDB, err := connectDatabase(cfg.Database.Local.Host, cfg.Database.Local.Port,
        cfg.Database.Local.User, cfg.Database.Local.Password, cfg.Database.Local.DBName,
        cfg.Database.Local.MaxConns, cfg.Database.Local.MinConns)
    if err != nil {
        log.Fatal("Failed to connect to local database:", err)
    }
    defer localDB.Close()

    // Initialize repositories
    mainRepo := mysql.NewMainRepository(mainDB, log)
    localRepo := mysql.NewLocalRepository(localDB, log)

    // Initialize domain service
    userService := user.NewDomainService(mainRepo, localRepo)

    // Initialize application handlers
    syncUserHandler := commands.NewSyncUserHandler(userService)
    syncAllHandler := commands.NewSyncAllUsersHandler(userService)
    getUserHandler := queries.NewGetUserHandler(userService)

    // Initialize HTTP handlers
    userHandler := handlers.NewUserHandler(
        syncUserHandler,
        syncAllHandler,
        getUserHandler,
        log,
        metrics,
    )

    // Setup router
    router := mux.NewRouter()
    
    // API routes
    api := router.PathPrefix("/api/v1").Subrouter()
    api.HandleFunc("/users/{id}/sync", userHandler.SyncUser).Methods("POST")
    api.HandleFunc("/users/sync/all", userHandler.SyncAllUsers).Methods("POST")
    api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")

    // Health check
    router.HandleFunc("/health", userHandler.HealthCheck).Methods("GET")

    // Metrics
    router.Handle("/metrics", promhttp.Handler())

    // Middleware
    router.Use(loggingMiddleware(log))
    router.Use(recoveryMiddleware(log))
    router.Use(corsMiddleware)

    // Create server
    server := &http.Server{
        Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
        Handler:      router,
        ReadTimeout:  cfg.Server.ReadTimeout,
        WriteTimeout: cfg.Server.WriteTimeout,
        IdleTimeout:  cfg.Server.IdleTimeout,
    }

    // Start server in goroutine
    go func() {
        log.Info("Starting server on port " + cfg.Server.Port)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal("Server failed:", err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    // Graceful shutdown
    log.Info("Shutting down server...")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Info("Server exited")
}

func connectDatabase(host string, port int, user, password, dbname string, maxConns, minConns int) (*sqlx.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=UTC&charset=utf8mb4",
        user, password, host, port, dbname)
    
    db, err := sqlx.Connect("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    db.SetMaxOpenConns(maxConns)
    db.SetMaxIdleConns(minConns)
    db.SetConnMaxLifetime(time.Hour)

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    return db, nil
}

func loggingMiddleware(log *zap.Logger) mux.MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            next.ServeHTTP(w, r)
            log.Info("HTTP request",
                zap.String("method", r.Method),
                zap.String("path", r.URL.Path),
                zap.Duration("duration", time.Since(start)),
                zap.String("remote_addr", r.RemoteAddr),
            )
        })
    }
}

func recoveryMiddleware(log *zap.Logger) mux.MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if err := recover(); err != nil {
                    log.Error("Panic recovered",
                        zap.Any("error", err),
                        zap.String("path", r.URL.Path),
                    )
                    http.Error(w, "Internal server error", http.StatusInternalServerError)
                }
            }()
            next.ServeHTTP(w, r)
        })
    }
}

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}