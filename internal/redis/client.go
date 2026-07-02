package redis

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/deba0208/stock-rsi-dashboard/internal/config"
	goredis "github.com/redis/go-redis/v9"
)

func NewClient(cfg *config.Config) (*goredis.Client, error) {

	client := goredis.NewClient(&goredis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			cfg.RedisHost,
			cfg.RedisPort,
		),
		Username: cfg.RedisUsername,
		Password: cfg.RedisPassword,
		TLSConfig: &tls.Config{},
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}