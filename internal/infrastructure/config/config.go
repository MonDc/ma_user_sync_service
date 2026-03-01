package config

import (
    "log"
    "time"
    
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    Observability ObservabilityConfig
}

type ServerConfig struct {
    Port         string
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    IdleTimeout  time.Duration
}

type DatabaseConfig struct {
    Main struct {
        Host     string
        Port     int
        User     string
        Password string
        DBName   string
        SSLMode  string
        MaxConns int
        MinConns int
    }
    Local struct {
        Host     string
        Port     int
        User     string
        Password string
        DBName   string
        SSLMode  string
        MaxConns int
        MinConns int
    }
}

type RedisConfig struct {
    Host     string
    Port     int
    Password string
    DB       int
}

type ObservabilityConfig struct {
    LogLevel     string
    MetricsPort  string
    Tracing      TracingConfig
}

type TracingConfig struct {
    Enabled     bool
    ServiceName string
    Environment string
    ExporterURL string
}

func LoadConfig() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./config")
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        log.Printf("Error reading config file: %s, using defaults", err)
    }

    config := &Config{
        Server: ServerConfig{
            Port:         viper.GetString("server.port"),
            ReadTimeout:  viper.GetDuration("server.read_timeout") * time.Second,
            WriteTimeout: viper.GetDuration("server.write_timeout") * time.Second,
            IdleTimeout:  viper.GetDuration("server.idle_timeout") * time.Second,
        },
    }

    // Database config
    config.Database.Main.Host = viper.GetString("database.main.host")
    config.Database.Main.Port = viper.GetInt("database.main.port")
    config.Database.Main.User = viper.GetString("database.main.user")
    config.Database.Main.Password = viper.GetString("database.main.password")
    config.Database.Main.DBName = viper.GetString("database.main.dbname")
    config.Database.Main.MaxConns = viper.GetInt("database.main.max_conns")
    config.Database.Main.MinConns = viper.GetInt("database.main.min_conns")

    config.Database.Local.Host = viper.GetString("database.local.host")
    config.Database.Local.Port = viper.GetInt("database.local.port")
    config.Database.Local.User = viper.GetString("database.local.user")
    config.Database.Local.Password = viper.GetString("database.local.password")
    config.Database.Local.DBName = viper.GetString("database.local.dbname")
    config.Database.Local.MaxConns = viper.GetInt("database.local.max_conns")
    config.Database.Local.MinConns = viper.GetInt("database.local.min_conns")

    // Observability
    config.Observability.LogLevel = viper.GetString("observability.log_level")
    config.Observability.MetricsPort = viper.GetString("observability.metrics_port")
    config.Observability.Tracing.Enabled = viper.GetBool("observability.tracing.enabled")
    config.Observability.Tracing.ServiceName = viper.GetString("observability.tracing.service_name")
    config.Observability.Tracing.Environment = viper.GetString("observability.tracing.environment")
    config.Observability.Tracing.ExporterURL = viper.GetString("observability.tracing.exporter_url")

    return config, nil
}