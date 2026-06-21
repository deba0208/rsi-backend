package repository

import (
	"context"
	"fmt"
	"strings"

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

	oldKeys, err := r.Client.Keys(ctx, "stock:*").Result()
	if err != nil && err != redis.Nil {
		return err
	}

	// Create a map of new symbols for quick lookup
	newStocksMap := make(map[string]bool)
	for _, stock := range stocks {
		newStocksMap[fmt.Sprintf("stock:%s", stock.Symbol)] = true
	}

	// Find stocks that are in Redis but no longer in the Nifty 50 list
	var keysToDelete []string
	var symbolsToDelete []string
	for _, oldKey := range oldKeys {
		if !newStocksMap[oldKey] {
			keysToDelete = append(keysToDelete, oldKey)
			symbolsToDelete = append(symbolsToDelete, strings.TrimPrefix(oldKey, "stock:"))
		}
	}

	// Delete the removed stocks, their cached metrics, and clear them from rankings
	if len(keysToDelete) > 0 {
		for _, symbol := range symbolsToDelete {
			keysToDelete = append(keysToDelete, fmt.Sprintf("metric:%s", symbol))
			r.Client.ZRem(ctx, "rsi:daily", symbol)
			r.Client.ZRem(ctx, "rsi:weekly", symbol)
			r.Client.ZRem(ctx, "rsi:monthly", symbol)
		}
		r.Client.Del(ctx, keysToDelete...)
	}

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
