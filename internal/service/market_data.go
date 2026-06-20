package service

import "github.com/deba0208/stock-rsi-dashboard/internal/models"

type MarketDataProvider interface {
	GetDailyCandles(symbol string) ([]models.Candle, error)

	GetWeeklyCandles(symbol string) ([]models.Candle, error)

	GetMonthlyCandles(symbol string) ([]models.Candle, error)
}