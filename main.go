package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/status", getStatus)
	router.GET("/transactions", getTransactions)
	router.GET("/transactions/:id", getTransaction)
	router.POST("/transactions", createTransaction)
	router.PUT("/transactions/:id", updateTransaction)
	router.DELETE("/transactions/:id", deleteTransaction)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Unable to start server. Error: ", err.Error())
	}
}
