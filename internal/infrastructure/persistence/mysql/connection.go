package mysql

import (
    "fmt"
    "time"

    _ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
    "go.uber.org/zap"
)

type Config struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    MaxConns int
    MinConns int
}

func NewConnection(cfg Config, logger *zap.Logger) (*sqlx.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=UTC&charset=utf8mb4",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
    
    db, err := sqlx.Connect("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(cfg.MaxConns)
    db.SetMaxIdleConns(cfg.MinConns)
    db.SetConnMaxLifetime(time.Hour)
    
    // Test connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    logger.Info("Connected to database",
        zap.String("host", cfg.Host),
        zap.String("database", cfg.DBName),
        zap.Int("max_conns", cfg.MaxConns),
    )
    
    return db, nil
}