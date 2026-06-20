package service

import (
	"time"

	"github.com/deba0208/stock-rsi-dashboard/internal/models"
)

type MockMarketDataProvider struct{}

func (m *MockMarketDataProvider) GetDailyCandles(
	symbol string,
) ([]models.Candle, error) {

	return []models.Candle{
		{Date: time.Now(), Close: 100},
		{Date: time.Now(), Close: 102},
		{Date: time.Now(), Close: 101},
		{Date: time.Now(), Close: 105},
		{Date: time.Now(), Close: 107},
		{Date: time.Now(), Close: 110},
		{Date: time.Now(), Close: 108},
		{Date: time.Now(), Close: 112},
		{Date: time.Now(), Close: 115},
		{Date: time.Now(), Close: 117},
		{Date: time.Now(), Close: 118},
		{Date: time.Now(), Close: 120},
		{Date: time.Now(), Close: 121},
		{Date: time.Now(), Close: 122},
		{Date: time.Now(), Close: 124},
	}, nil
}

func (m *MockMarketDataProvider) GetWeeklyCandles(
	symbol string,
) ([]models.Candle, error) {
	return []models.Candle{
		{Date: time.Now(), Close: 100},
		{Date: time.Now(), Close: 102},
		{Date: time.Now(), Close: 101},
		{Date: time.Now(), Close: 105},
		{Date: time.Now(), Close: 107},
		{Date: time.Now(), Close: 110},
		{Date: time.Now(), Close: 108},
		{Date: time.Now(), Close: 112},
		{Date: time.Now(), Close: 115},
		{Date: time.Now(), Close: 117},
		{Date: time.Now(), Close: 118},
		{Date: time.Now(), Close: 120},
		{Date: time.Now(), Close: 121},
		{Date: time.Now(), Close: 122},
		{Date: time.Now(), Close: 124},
	}, nil
}

func (m *MockMarketDataProvider) GetMonthlyCandles(
	symbol string,
) ([]models.Candle, error) {
	return []models.Candle{
		{Date: time.Now(), Close: 100},
		{Date: time.Now(), Close: 102},
		{Date: time.Now(), Close: 101},
		{Date: time.Now(), Close: 105},
		{Date: time.Now(), Close: 107},
		{Date: time.Now(), Close: 110},
		{Date: time.Now(), Close: 108},
		{Date: time.Now(), Close: 112},
		{Date: time.Now(), Close: 115},
		{Date: time.Now(), Close: 117},
		{Date: time.Now(), Close: 118},
		{Date: time.Now(), Close: 120},
		{Date: time.Now(), Close: 121},
		{Date: time.Now(), Close: 122},
		{Date: time.Now(), Close: 124},
	}, nil
}
