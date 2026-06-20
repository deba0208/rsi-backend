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
	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	return &Config{
		Port:          os.Getenv("PORT"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}, nil
}
