package repository

import (
	"context"
	"fmt"

	"github.com/deba0208/stock-rsi-dashboard/internal/models"
	"github.com/redis/go-redis/v9"
)

type StockRepository struct {
	Client *redis.Client
}

func NewStockRepository(client *redis.Client) *StockRepository {
	return &StockRepository{Client: client}
}

func (r *StockRepository) SaveStocks(stocks []models.Stock) error {
	ctx := context.Background()

	for _, stock := range stocks {
		key := fmt.Sprintf("stock:%s", stock.Symbol)

		// Store each field as a named hash field
		if err := r.Client.HSet(ctx, key,
			"symbol", stock.Symbol,
			"name", stock.Name,
		).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (r *StockRepository) GetAllStocks() ([]models.Stock, error) {
	ctx := context.Background()

	keys, err := r.Client.Keys(ctx, "stock:*").Result()
	if err != nil {
		return nil, err
	}

	var stocks []models.Stock
	for _, key := range keys {
		// HGetAll returns map[string]string of all hash fields
		fields, err := r.Client.HGetAll(ctx, key).Result()
		if err != nil || len(fields) == 0 {
			continue
		}

		stocks = append(stocks, models.Stock{
			Symbol: fields["symbol"],
			Name:   fields["name"],
		})
	}

	return stocks, nil
}
