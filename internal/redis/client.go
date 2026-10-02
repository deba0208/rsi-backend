
package redis

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/deba0208/stock-rsi-dashboard/internal/config"
    goredis "github.com/redis/go-redis/v9"
)

func NewClient(cfg *config.Config) (*goredis.Client, error) {
    addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

    client := goredis.NewClient(&goredis.Options{
        Addr:         addr,
        Username:     cfg.RedisUsername,
        Password:     cfg.RedisPassword,
        DialTimeout:  5 * time.Second,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 5 * time.Second,
        MaxRetries:   0,
        Protocol:     2,
        OnConnect: func(ctx context.Context, cn *goredis.Conn) error {
            log.Println("Redis TCP connection established")
            return nil
        },
    })

    ctx, cancel := context.WithTimeout(
        context.Background(),
        10*time.Second,
    )
    defer cancel()

    log.Printf("Sending Redis PING to %s", addr)

    if err := client.Ping(ctx).Err(); err != nil {
        _ = client.Close()
        return nil, fmt.Errorf("Redis PING failed: %w", err)
    }

    log.Println("Redis PING successful")
    return client, nil
}