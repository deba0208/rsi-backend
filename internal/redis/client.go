package redis

import (
	"context"
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
		Password: cfg.RedisPassword,
	})

	err := client.Ping(
		context.Background(),
	).Err()

	if err != nil {
		return nil, err
	}

	return client, nil
}
