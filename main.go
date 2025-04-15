package main

import (
	"log"
	"strings"

	"slices"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.Use(corsMiddleware())

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

func corsMiddleware() gin.HandlerFunc {
	originString := "http://localhost:4200, https://tavo826.github.io/, https://transcendent-brioche-97eea6.netlify.app"
	var allowedOrigins []string
	if originString != "" {
		allowedOrigins = strings.Split(originString, ",")
	}

	return func(ctx *gin.Context) {
		isOriginAllowed := func(origin string, allowedOrigins []string) bool {
			return slices.Contains(allowedOrigins, origin)
		}

		origin := ctx.Request.Header.Get("Origin")

		if isOriginAllowed(origin, allowedOrigins) {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, GET, PUT")
		}

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}
