package main

import (
	"log"

	"github.com/deba0208/stock-rsi-dashboard/internal/config"
	"github.com/deba0208/stock-rsi-dashboard/internal/handler"
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
	nifty50Provider := service.NewNifty50Provider()
	stockService := service.NewStockService(stockRepos, nifty50Provider)

	metricRepo := repository.NewMetricRepository(client)
	metricService := service.NewMetricService(rsiService, marketProvider, metricRepo)

	// --- Scheduler ---
	rsiScheduler := scheduler.NewRSIScheduler(stockService, metricService)
	scheduler.Start(rsiScheduler)
	// --- Router ---
	router := gin.Default()

	// Health
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// Top 50 by criteria (daily / weekly / monthly)
	metricHandler := handler.NewMetricHandler(metricService)
	router.GET("/metrics/top50", metricHandler.GetTop50ByCriteria)

	// router.Run() blocks — must be the very last call
	router.Run(":" + cfg.Port)
}
