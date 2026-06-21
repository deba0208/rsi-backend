package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/deba0208/stock-rsi-dashboard/internal/config"
	"github.com/deba0208/stock-rsi-dashboard/internal/handler"
	"github.com/deba0208/stock-rsi-dashboard/internal/models"
	"github.com/deba0208/stock-rsi-dashboard/internal/redis"
	"github.com/deba0208/stock-rsi-dashboard/internal/repository"
	"github.com/deba0208/stock-rsi-dashboard/internal/scheduler"
	"github.com/deba0208/stock-rsi-dashboard/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	client, err := redis.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// --- Services ---
	marketProvider := service.NewYahooMarketDataService()
	rsiService := service.NewRSIService(marketProvider)

	stockRepos := repository.NewStockRepository(client)
	stockService := service.NewStockService(stockRepos)

	metricRepo := repository.NewMetricRepository(client)
	metricService := service.NewMetricService(rsiService, marketProvider, metricRepo)

	// Initialize stocks from JSON
	// NSE_STOCKS_PATH allows overriding the stocks file path at runtime;
	// defaults to a path relative to where the binary is invoked.
	stocksPath := os.Getenv("NSE_STOCKS_PATH")
	if stocksPath == "" {
		stocksPath = "./internal/config/nse_stocks.json"
	}
	err = stockService.InitializeStocks(stocksPath)
	if err != nil {
		fmt.Println("Warning: could not initialize stocks:", err)
	}

	// --- Scheduler ---
	rsiScheduler := scheduler.NewRSIScheduler(stockService, metricService)

	// --- Router ---
	router := gin.Default()

	// Health
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// RSI by query params
	router.GET("/rsi", func(c *gin.Context) {
		symbol := c.Query("symbol")
		timeFrame := c.Query("timeFrame")

		if symbol == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing required query param: symbol"})
			return
		}
		if timeFrame == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing required query param: timeFrame"})
			return
		}

		rsiValue, err := rsiService.RSI(symbol, timeFrame)
		if err != nil {
			// Invalid timeframe is a client error; everything else is a server error
			status := http.StatusInternalServerError
			if err.Error() == fmt.Sprintf("invalid timeframe: %s", timeFrame) {
				status = http.StatusBadRequest
			}
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"symbol":    symbol,
			"timeFrame": timeFrame,
			"rsi":       rsiValue,
		})
	})

	// Manual Scheduler Trigger
	router.POST(
		"/scheduler/run",
		func(c *gin.Context) {
			// Run async so the HTTP response returns immediately.
			// Progress and errors are logged server-side.
			go rsiScheduler.Run()
			c.JSON(202, gin.H{"message": "scheduler started"})
		},
	)

	// Top 50 by criteria (daily / weekly / monthly)
	metricHandler := handler.NewMetricHandler(metricService)
	router.GET("/metrics/top50/:criteria", metricHandler.GetTop50ByCriteria)

	// Yahoo Finance candles endpoint
	// Reuse marketProvider — avoids a second http.Client allocation
	provider := marketProvider

	router.GET("/candles", func(c *gin.Context) {

		symbol := c.Query("symbol")
		timeFrame := c.Query("timeFrame")

		var (
			candles []models.Candle
			err     error
		)

		switch timeFrame {

		case "daily":
			candles, err =
				provider.GetCandles(symbol, "1d")

		case "weekly":
			candles, err =
				provider.GetCandles(symbol, "1wk")

		case "monthly":
			candles, err =
				provider.GetCandles(symbol, "1mo")

		default:
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": "invalid timeframe",
				},
			)
			return
		}

		if err != nil {

			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"error": err.Error(),
				},
			)

			return
		}

		c.JSON(
			http.StatusOK,
			candles,
		)
	})

	log.Println("Server is running on :" + cfg.Port)

	// router.Run() blocks — must be the very last call
	router.Run(":" + cfg.Port)
}
