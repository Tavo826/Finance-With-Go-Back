package http

import (
	"personal-finance/adapter/config"
	"personal-finance/adapter/handler/http/token"
	"strings"
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
	authHandler AuthHandler,
) (*Router, error) {

	if config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	middleware := token.NewAuthMiddleware()

	allowedOrigins := strings.Split(config.App.AllowedOrigins, ",")

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
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

		auth := v1.Group("/users")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
		auth.Use(middleware.Implement(config.Token))
		{
			auth.GET("/", authHandler.GetUserById)
			auth.PUT("/:id", authHandler.UpdateUser)
			auth.DELETE("/:id", authHandler.DeleteUser)
		}

		transaction := v1.Group("/transactions")
		transaction.Use(middleware.Implement(config.Token))
		{
			transaction.GET("/", transactionHandler.GetTransactionsByUserId)
			transaction.GET("/filter_date", transactionHandler.GetTransactionsByDate)
			transaction.GET("/filter_subject", transactionHandler.GetTransactionsBySubject)
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
