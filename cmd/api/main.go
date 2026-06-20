package main

import (
	"fmt"
	"log"

	"github.com/deba0208/stock-rsi-dashboard/internal/config"
	"github.com/deba0208/stock-rsi-dashboard/internal/handler"
	"github.com/deba0208/stock-rsi-dashboard/internal/redis"
	"github.com/deba0208/stock-rsi-dashboard/internal/repository"
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
	mockProvider := &service.MockMarketDataProvider{}
	rsiService := service.NewRSIService(mockProvider)

	stockRepos := repository.NewStockRepository(client)
	stockService := service.NewStockService(stockRepos)

	metricRepo := repository.NewMetricRepository(client)
	metricService := service.NewMetricService(rsiService, metricRepo)

	// Initialize stocks from JSON
	err = stockService.InitializeStocks("./internal/config/nse_stocks.json")
	if err != nil {
		fmt.Println("Warning: could not initialize stocks:", err)
	}

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
		rsiValue, err := rsiService.RSI(symbol, timeFrame)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"symbol":    symbol,
			"timeFrame": timeFrame,
			"rsi":       rsiValue,
		})
	})

	// Top 50 by criteria (daily / weekly / monthly)
	metricHandler := handler.NewMetricHandler(metricService)
	router.GET("/metrics/top50/:criteria", metricHandler.GetTop50ByCriteria)

	log.Println("Server is running on :" + cfg.Port)

	// router.Run() blocks — must be last
	router.Run(":" + cfg.Port)
}
