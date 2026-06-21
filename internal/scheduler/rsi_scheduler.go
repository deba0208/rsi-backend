package scheduler

import (
	"log"

	"github.com/deba0208/stock-rsi-dashboard/internal/service"
)

type RSIScheduler struct {
	stockService  *service.StockService
	metricService *service.MetricService
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

func (s *RSIScheduler) Run() {
	stocks, err := s.stockService.GetStocks()
	if err != nil {
		log.Printf("[scheduler] failed to fetch stocks: %v", err)
		return
	}

	log.Printf("[scheduler] starting RSI update for %d stocks", len(stocks))

	for _, stock := range stocks {

		log.Printf(
			"Updating %s",
			stock.Symbol,
		)

		err :=
			s.metricService.
				UpdateMetric(
					stock.Symbol,
				)

		if err != nil {

			log.Printf(
				"Failed %s : %v",
				stock.Symbol,
				err,
			)

			continue
		}
	}

	log.Printf("[scheduler] RSI update complete")
}
