package http

import (
	"personal-finance/adapter/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine
}

func NewRouter(
	config *config.Container,
	transactionHandler TransactionHandler,
) (*Router, error) {

	if config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200", "https://tavo826.github.io", "https://transcendent-brioche-97eea6.netlify.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := router.Group("/v1")
	{
		status := v1.Group("")
		{
			status.GET("/", transactionHandler.GetStatus)
			status.GET("/status", transactionHandler.GetStatus)
		}

		transaction := v1.Group("/transactions")
		{
			transaction.GET("/", transactionHandler.GetTransactions)
			transaction.GET("/:id", transactionHandler.GetTransaction)
			transaction.POST("/", transactionHandler.CreateTransaction)
			transaction.PUT("/:id", transactionHandler.UpdateTransaction)
			transaction.DELETE("/:id", transactionHandler.DeleteTransaction)
		}
	}

	return &Router{
		router,
	}, nil
}

func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
