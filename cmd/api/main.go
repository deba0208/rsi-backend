package main

import (
	"fmt"
	"log"

	"github.com/deba0208/stock-rsi-dashboard/internal/config"
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
		fmt.Println(err)
	}

	fmt.Println(client)

	router := gin.Default()

	// Services — must be declared before route handlers that use them
	mockProvider := &service.MockMarketDataProvider{}
	rsiService := service.NewRSIService(mockProvider)

	stockRepos := repository.NewStockRepository(client)
	stockService := service.NewStockService(stockRepos)

	err = stockService.InitializeStocks("./internal/config/nse_stocks.json")
	if err != nil {
		fmt.Println(err)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})

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

	log.Println("Server is running on :" + cfg.Port)

	router.Run(":" + cfg.Port)
}
