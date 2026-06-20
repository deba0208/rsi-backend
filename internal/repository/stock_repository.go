package repository

import (
	"context"
	"encoding/json"
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

		data, err := json.Marshal(stock)
		if err != nil {
			return err
		}

		if err := r.Client.Set(ctx, key, data, 0).Err(); err != nil {
			return err
		}
	}

	return nil
}