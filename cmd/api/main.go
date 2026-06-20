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

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})

	stockRepos := repository.NewStockRepository(client)
	
	stockService := service.NewStockService(stockRepos)

	err = stockService.InitializeStocks("./internal/config/nse_stocks.json")
	if err != nil {
		fmt.Println(err)
	}

	log.Println("Server is running on :" + cfg.Port)

	router.Run(":" + cfg.Port)
}
