package service

import (
	"fmt"

	"github.com/deba0208/stock-rsi-dashboard/internal/models"
	"github.com/deba0208/stock-rsi-dashboard/internal/rsi"
)

// Timeframe constants
const (
	Daily   = "daily"
	Weekly  = "weekly"
	Monthly = "monthly"
)

type RSIService struct {
	provider MarketDataProvider
}

func NewRSIService(
	provider MarketDataProvider,
) *RSIService {

	return &RSIService{
		provider: provider,
	}
}

func (s *RSIService) RSI(
	symbol string,
	timeFrame string,
) (float64, error) {

	var candles []models.Candle
	var err error

	switch timeFrame {
	case Daily:
		candles, err = s.provider.GetDailyCandles(symbol)
	case Weekly:
		candles, err = s.provider.GetWeeklyCandles(symbol)
	case Monthly:
		candles, err = s.provider.GetMonthlyCandles(symbol)
	default:
		return 0, fmt.Errorf("invalid timeframe: %s", timeFrame)
	}

	// Bug fix: check error after fetching candles
	if err != nil {
		return 0, fmt.Errorf("failed to fetch candles for %s (%s): %w", symbol, timeFrame, err)
	}

	if len(candles) == 0 {
		return 0, fmt.Errorf("no candle data for %s (%s)", symbol, timeFrame)
	}

	closes := make([]float64, 0, len(candles))

	for _, candle := range candles {
		closes = append(closes, candle.Close)
	}

	return rsi.Calculate(closes, 14), nil
}