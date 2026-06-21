package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	RedisHost     string
	RedisPort     string
	RedisPassword string
}

func LoadConfig() (*Config, error) {
	// Ignore error — .env file may not exist in production/container environments
	// where env vars are injected directly.
	_ = godotenv.Load()

	return &Config{
		Port:          os.Getenv("PORT"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}, nil
}
