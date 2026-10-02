
package config

import (
    "fmt"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Port          string
    RedisHost     string
    RedisPort     string
    RedisUsername string
    RedisPassword string
}

func LoadConfig() (*Config, error) {
    // .env is optional in production
    _ = godotenv.Load()

    cfg := &Config{
        Port:          os.Getenv("PORT"),
        RedisHost:     os.Getenv("REDIS_HOST"),
        RedisPort:     os.Getenv("REDIS_PORT"),
        RedisUsername: os.Getenv("REDIS_USERNAME"),
        RedisPassword: os.Getenv("REDIS_PASSWORD"),
    }

    if cfg.Port == "" {
        cfg.Port = "8080"
    }

    if cfg.RedisHost == "" || cfg.RedisPort == "" {
        return nil, fmt.Errorf("REDIS_HOST and REDIS_PORT are required")
    }

    if cfg.RedisPassword == "" {
        return nil, fmt.Errorf("REDIS_PASSWORD is required")
    }

    return cfg, nil
}