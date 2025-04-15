package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200", "https://tavo826.github.io", "https://transcendent-brioche-97eea6.netlify.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/", getStatus)
	router.GET("/status", getStatus)
	router.GET("/transactions", getTransactions)
	router.GET("/transactions/:id", getTransaction)
	router.POST("/transactions", createTransaction)
	router.PUT("/transactions/:id", updateTransaction)
	router.DELETE("/transactions/:id", deleteTransaction)

	if err := router.Run(":8000"); err != nil {
		log.Fatal("Unable to start server. Error: ", err.Error())
	}
}
