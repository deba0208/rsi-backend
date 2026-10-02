
package scheduler

import (
	"fmt"
	"log"
	"sync"

	"github.com/deba0208/stock-rsi-dashboard/internal/service"
)

type RSIScheduler struct {
	stockService  *service.StockService
	metricService *service.MetricService

	mu      sync.Mutex
	running bool
}

func NewRSIScheduler(
	stockService *service.StockService,
	metricService *service.MetricService,
) *RSIScheduler {
	return &RSIScheduler{
		stockService:  stockService,
		metricService: metricService,
	}
}

func (r *RSIScheduler) Run() error {
	// Prevent concurrent executions.
	r.mu.Lock()
	if r.running {
		r.mu.Unlock()
		log.Println("RSI update already running; skipping this execution")
		return nil
	}
	r.running = true
	r.mu.Unlock()

	defer func() {
		r.mu.Lock()
		r.running = false
		r.mu.Unlock()
	}()

	log.Println("Starting RSI update...")

	// Initialize or refresh the stock list.
	if err := r.stockService.InitializeStocks(); err != nil {
		log.Printf("Warning: failed to initialize stocks: %v", err)
	}

	// Fetch the stocks to process.
	stocks, err := r.stockService.GetStocks()
	if err != nil {
		return fmt.Errorf("failed to get stocks: %w", err)
	}

	if len(stocks) == 0 {
		log.Println("No stocks found for RSI update")
		return nil
	}

	log.Printf("Updating RSI for %d stocks", len(stocks))

	failed := 0

	// Process stocks sequentially.
	for _, stock := range stocks {
		if err := r.metricService.UpdateMetric(stock.Symbol); err != nil {
			failed++
			log.Printf(
				"Failed to update RSI for %s: %v",
				stock.Symbol,
				err,
			)
			continue
		}

		log.Printf("Successfully updated RSI for %s", stock.Symbol)
	}

	log.Printf(
		"RSI update completed: %d stocks processed, %d failed",
		len(stocks),
		failed,
	)

	if failed > 0 {
		return fmt.Errorf("%d of %d stock updates failed", failed, len(stocks))
	}

	return nil
}