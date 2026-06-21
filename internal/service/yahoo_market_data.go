package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/deba0208/stock-rsi-dashboard/internal/models"
)

type YahooMarketDataService struct {
	client *http.Client
}

func NewYahooMarketDataService() *YahooMarketDataService {
	return &YahooMarketDataService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

const (
	intervalDaily   = "1d"
	intervalWeekly  = "1wk"
	intervalMonthly = "1mo"
)

func (y *YahooMarketDataService) fetchCandles(symbol, rangeVal,
	interval string) ([]models.Candle, error) {

	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s.NS?range=%s&interval=%s", symbol, rangeVal, interval)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	// Yahoo Finance blocks requests without a browser-like User-Agent (returns 429)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := y.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch yahoo data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yahoo returned non-OK status: %d", resp.StatusCode)
	}

	var yahooResp models.YahooResponse

	err = json.NewDecoder(resp.Body).
		Decode(&yahooResp)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to decode yahoo response: %w",
			err,
		)
	}

	if len(yahooResp.Chart.Result) == 0 {
		return nil, fmt.Errorf(
			"no data found for symbol %s",
			symbol,
		)
	}

	result := yahooResp.Chart.Result[0]

	if len(result.Indicators.Quote) == 0 {
		return nil, fmt.Errorf(
			"quote data missing for symbol %s",
			symbol,
		)
	}

	candles := make(
		[]models.Candle,
		0,
		len(result.Timestamp),
	)

	closePrices :=
		result.Indicators.Quote[0].Close

	for i, ts := range result.Timestamp {

		if i >= len(closePrices) {
			break
		}

		if closePrices[i] == nil {
			continue
		}

		candles = append(
			candles,
			models.Candle{
				Date:  time.Unix(ts, 0),
				Close: *closePrices[i],
			},
		)
	}

	if len(candles) == 0 {
		return nil, fmt.Errorf(
			"no valid candles found for symbol %s",
			symbol,
		)
	}

	return candles, nil
}

func (y *YahooMarketDataService) GetCandles(symbol string,
	interval string) ([]models.Candle, error) {
	var rangeVal string
	switch interval {
	case intervalDaily:
		rangeVal = "2y"
	case intervalWeekly:
		rangeVal = "10y"
	case intervalMonthly:
		rangeVal = "20y"
	default:
		return nil, fmt.Errorf("invalid interval: %s", interval)
	}

	return y.fetchCandles(symbol, rangeVal, interval)
}

// GetCurrentPrice returns the most recent close price for a symbol.
// It fetches the last 5 daily candles and returns the last one.
func (y *YahooMarketDataService) GetCurrentPrice(symbol string) (float64, error) {
	candles, err := y.fetchCandles(symbol, "5d", intervalDaily)
	if err != nil {
		return 0, fmt.Errorf("failed to get current price for %s: %w", symbol, err)
	}

	// Return the most recent close (last element)
	return candles[len(candles)-1].Close, nil
}
