package main

import (
	"log"
	"time"

	"github.com/deba0208/stock-rsi-dashboard/internal/config"
	"github.com/deba0208/stock-rsi-dashboard/internal/handler"
	"github.com/deba0208/stock-rsi-dashboard/internal/redis"
	"github.com/deba0208/stock-rsi-dashboard/internal/repository"
	"github.com/deba0208/stock-rsi-dashboard/internal/scheduler"
	"github.com/deba0208/stock-rsi-dashboard/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	log.Println("Loading application configuration")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Redis address: %s:%s", cfg.RedisHost, cfg.RedisPort)
	log.Printf("Redis username: %s", cfg.RedisUsername)
	log.Println("Initializing Redis client...")

	client, err := redis.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Redis connected successfully")
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
	if err := scheduler.Start(rsiScheduler); err != nil {
		log.Fatalf("Failed to start RSI scheduler: %v", err)
	}
	// --- Router ---
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5174",
			"https://rsi-frontend.vercel.app/",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

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
