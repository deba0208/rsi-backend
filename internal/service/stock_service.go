package service

import (
	"encoding/json"
	"log"
	"os"

	"github.com/deba0208/stock-rsi-dashboard/internal/models"
	"github.com/deba0208/stock-rsi-dashboard/internal/repository"
)

type StockService struct {
	Repo            *repository.StockRepository
	Nifty50Provider *Nifty50Provider
}

func NewStockService(
	repo *repository.StockRepository,
	nifty50Provider *Nifty50Provider,
) *StockService {
	return &StockService{
		Repo:            repo,
		Nifty50Provider: nifty50Provider,
	}
}

// InitializeStocks seeds the stock list from the live Nifty50 CSV feed.
// If the feed is unavailable, it falls back to the local JSON file.
func (s *StockService) InitializeStocks() error {
	stocks, err := s.Nifty50Provider.GetStocks()
	if err != nil {
		log.Printf("[stocks] live feed unavailable (%v), falling back to local JSON", err)
		stocks, err = loadStocksFromJSON("./internal/config/nse_stocks.json")
		if err != nil {
			return err
		}
	}
	return s.Repo.SaveStocks(stocks)
}

func (s *StockService) GetStocks() ([]models.Stock, error) {
	return s.Repo.GetAllStocks()
}

// loadStocksFromJSON reads a JSON file and returns the stock list.
func loadStocksFromJSON(filePath string) ([]models.Stock, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var stocks []models.Stock
	if err := json.Unmarshal(data, &stocks); err != nil {
		return nil, err
	}
	return stocks, nil
}
