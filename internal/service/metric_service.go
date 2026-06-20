package service

import (
	"github.com/deba0208/stock-rsi-dashboard/internal/models"
	"github.com/deba0208/stock-rsi-dashboard/internal/repository"
)

type MetricService struct {
	rsiService *RSIService
	repo       *repository.MetricRepository
}

func NewMetricService(
	rsiService *RSIService,
	repo *repository.MetricRepository,
) *MetricService {

	return &MetricService{
		rsiService: rsiService,
		repo:       repo,
	}
}

func (s *MetricService) UpdateMetric(
	symbol string,
) error {

	dailyRSI, err :=
		s.rsiService.RSI(symbol, "daily")

	if err != nil {
		return err
	}

	weeklyRSI, err :=
		s.rsiService.RSI(symbol, "weekly")

	if err != nil {
		return err
	}

	monthlyRSI, err :=
		s.rsiService.RSI(symbol, "monthly")

	if err != nil {
		return err
	}

	metric := models.StockMetric{
		Symbol:     symbol,
		Price:      0, // later from market feed
		DailyRSI:   dailyRSI,
		WeeklyRSI:  weeklyRSI,
		MonthlyRSI: monthlyRSI,
	}

	return s.repo.SaveMetric(metric)
}

func (s *MetricService) GetTop50ByCriteria(criteria string) ([]string, error) {
	return s.repo.GetTop50ByCriteria(criteria)
}
