package service

import (
	"encoding/json"
	"os"

	"github.com/deba0208/stock-rsi-dashboard/internal/models"
)

func LoadStocks(filePath string) ([]models.Stock, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var stocks []models.Stock
	err = json.Unmarshal(data, &stocks)
	if err != nil {
		return nil, err
	}

	return stocks, nil
}