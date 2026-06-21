package service

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"github.com/deba0208/stock-rsi-dashboard/internal/models"
)

const nifty50URL = "https://www.niftyindices.com/IndexConstituent/ind_nifty50list.csv"

type Nifty50Provider struct {
	client *http.Client
}

func NewNifty50Provider() *Nifty50Provider {
	return &Nifty50Provider{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (p *Nifty50Provider) GetStocks() ([]models.Stock, error) {

	req, err := http.NewRequest(http.MethodGet, nifty50URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	// niftyindices.com blocks requests without browser-like headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.niftyindices.com/")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch nifty50 list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("niftyindices.com returned status: %d", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	stocks := make([]models.Stock, 0, 50)

	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}
		if len(row) < 3 {
			continue
		}
		stocks = append(stocks, models.Stock{
			Symbol: row[2],
			Name:   row[1],
		})
	}

	if len(stocks) == 0 {
		return nil, fmt.Errorf("no stocks found in nifty50 list")
	}

	return stocks, nil
}
