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

    "go.uber.org/zap"
    "github.com/spf13/viper"
    "github.com/jmoiron/sqlx"
    "github.com/mondc/ma_user_sync_service/internal/infrastructure/persistence/mysql"
)

func initConfig() error {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./config")
    viper.AddConfigPath(".")
    viper.AutomaticEnv()

    // Default values
    viper.SetDefault("server.port", "8080")
    viper.SetDefault("server.read_timeout", 15)
    viper.SetDefault("server.write_timeout", 15)
    viper.SetDefault("server.idle_timeout", 60)
    viper.SetDefault("log.level", "info")

    if err := viper.ReadInConfig(); err != nil {
        // It's okay if config file doesn't exist, we use defaults
        log.Printf("No config file found, using defaults: %v", err)
    }

    return nil
}

func initLogger() (*zap.Logger, error) {
    logLevel := viper.GetString("log.level")
    
    var config zap.Config
    if logLevel == "development" {
        config = zap.NewDevelopmentConfig()
    } else {
        config = zap.NewProductionConfig()
        config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
    }
    
    return config.Build()
}

func setupRouter(logger *zap.Logger) *http.ServeMux {
    router := http.NewServeMux()
    
    // Health check endpoint
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"healthy","service":"ma-user-sync"}`))
    })
    
    // API routes will be added here
    router.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message":"Users endpoint placeholder"}`))
    })
    
    return router
}

func main() {
    // 1. Load configuration
    if err := initConfig(); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // 2. Initialize logger
    logger, err := initLogger()
    if err != nil {
        log.Fatalf("Failed to init logger: %v", err)
    }
    defer logger.Sync()

    logger.Info("Starting MA User Sync Service")

    // 3. Setup HTTP server
    router := setupRouter(logger)
    server := &http.Server{
        Addr:         fmt.Sprintf(":%s", viper.GetString("server.port")),
        Handler:      router,
        ReadTimeout:  viper.GetDuration("server.read_timeout") * time.Second,
        WriteTimeout: viper.GetDuration("server.write_timeout") * time.Second,
        IdleTimeout:  viper.GetDuration("server.idle_timeout") * time.Second,
    }

    // 3.2 Initialize databases
    miDB, maDB, err := initDatabases(logger)
    if err != nil {
        logger.Fatal("Failed to connect to databases", zap.Error(err))
    }
    defer miDB.Close()
    defer maDB.Close()

    logger.Info("Both databases connected successfully")
    

    // 4. Start server in goroutine
    go func() {
        logger.Info("Server listening", zap.String("port", viper.GetString("server.port")))
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatal("Server failed", zap.Error(err))
        }
    }()

    // 5. Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    logger.Info("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        logger.Fatal("Server forced to shutdown", zap.Error(err))
    }

    logger.Info("Server exited")
}

func initDatabases(logger *zap.Logger) (*sqlx.DB, *sqlx.DB, error) {
    // MI Database (main users)
    miDB, err := mysql.NewConnection(mysql.Config{
        Host:     viper.GetString("database.mi.host"),
        Port:     viper.GetInt("database.mi.port"),
        User:     viper.GetString("database.mi.user"),
        Password: viper.GetString("database.mi.password"),
        DBName:   viper.GetString("database.mi.dbname"),
        MaxConns: viper.GetInt("database.mi.max_conns"),
        MinConns: viper.GetInt("database.mi.min_conns"),
    }, logger.With(zap.String("db", "mi")))
    if err != nil {
        return nil, nil, fmt.Errorf("MI database: %w", err)
    }
    
    // MA Database (local users)
    maDB, err := mysql.NewConnection(mysql.Config{
        Host:     viper.GetString("database.ma.host"),
        Port:     viper.GetInt("database.ma.port"),
        User:     viper.GetString("database.ma.user"),
        Password: viper.GetString("database.ma.password"),
        DBName:   viper.GetString("database.ma.dbname"),
        MaxConns: viper.GetInt("database.ma.max_conns"),
        MinConns: viper.GetInt("database.ma.min_conns"),
    }, logger.With(zap.String("db", "ma")))
    if err != nil {
        return nil, nil, fmt.Errorf("MA database: %w", err)
    }
    
    return miDB, maDB, nil
}