package service

import "github.com/deba0208/stock-rsi-dashboard/internal/models"

type MarketDataProvider interface {
	GetCandles(symbol, timeframe string) ([]models.Candle, error)
	GetCurrentPrice(symbol string) (float64, error)
}
