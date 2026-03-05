package config

import (
    "github.com/spf13/viper"
    "github.com/joho/godotenv"
)

func Load() error {
    // Load .env file first (if exists)
    godotenv.Load("config/.env")
    
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./config")
    
    // Enable env var substitution
    viper.AutomaticEnv()
    
    return viper.ReadInConfig()
}