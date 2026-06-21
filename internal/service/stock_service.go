package service

import (
	"github.com/deba0208/stock-rsi-dashboard/internal/models"
	"github.com/deba0208/stock-rsi-dashboard/internal/repository"
)

type StockService struct {
	Repo *repository.StockRepository
}

func NewStockService(repo *repository.StockRepository) *StockService {
	return &StockService{Repo: repo}
}

func (s *StockService) InitializeStocks(filePath string) error {
	stocks, err := LoadStocks(filePath)
	if err != nil {
		return err
	}
	return s.Repo.SaveStocks(stocks)
}

func (s *StockService) GetStocks() ([]models.Stock, error) {
	return s.Repo.GetAllStocks()
}
