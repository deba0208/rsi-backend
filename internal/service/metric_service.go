package service

import (
	"log"

	"github.com/deba0208/stock-rsi-dashboard/internal/models"
	"github.com/deba0208/stock-rsi-dashboard/internal/repository"
)

type MetricService struct {
	rsiService     *RSIService
	marketProvider *YahooMarketDataService
	repo           *repository.MetricRepository
}

func NewMetricService(
	rsiService *RSIService,
	marketProvider *YahooMarketDataService,
	repo *repository.MetricRepository,
) *MetricService {

	return &MetricService{
		rsiService:     rsiService,
		marketProvider: marketProvider,
		repo:           repo,
	}
}

func (s *MetricService) UpdateMetric(
	symbol string,
) error {

	// Fetch each timeframe independently — a failure in one does not
	// discard the others. The metric is always saved with whatever
	// data was available (failed timeframes are stored as 0).
	dailyRSI, err := s.rsiService.RSI(symbol, "daily")
	if err != nil {
		log.Printf("[metric] %s daily RSI failed: %v", symbol, err)
		dailyRSI = 0
	}

	weeklyRSI, err := s.rsiService.RSI(symbol, "weekly")
	if err != nil {
		log.Printf("[metric] %s weekly RSI failed: %v", symbol, err)
		weeklyRSI = 0
	}

	monthlyRSI, err := s.rsiService.RSI(symbol, "monthly")
	if err != nil {
		log.Printf("[metric] %s monthly RSI failed: %v", symbol, err)
		monthlyRSI = 0
	}

	price, err := s.marketProvider.GetCurrentPrice(symbol)
	if err != nil {
		log.Printf("[metric] %s current price failed: %v", symbol, err)
		price = 0
	}

	metric := models.StockMetric{
		Symbol:     symbol,
		Price:      price,
		DailyRSI:   dailyRSI,
		WeeklyRSI:  weeklyRSI,
		MonthlyRSI: monthlyRSI,
	}

	return s.repo.SaveMetric(metric)
}

func (s *MetricService) GetTop50ByCriteria(criteria string) ([]string, error) {
	return s.repo.GetTop50ByCriteria(criteria)
}
