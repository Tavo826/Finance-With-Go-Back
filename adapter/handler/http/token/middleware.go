package token

import (
	"personal-finance/adapter/config"
	"personal-finance/adapter/handler/http/dto"
	"personal-finance/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct{}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (am *AuthMiddleware) Implement(config *config.Token) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		jwtSecret := []byte(config.JwtSecret)

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			dto.HandleError(ctx, domain.ErrInvalidAuthorizationHeader)
			ctx.Abort()
			return
		}

		tokenString := authHeader[7:]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			dto.HandleError(ctx, domain.ErrInvalidToken)
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			dto.HandleError(ctx, domain.ErrTokenCreation)
			ctx.Abort()
			return
		}

		userID := claims["id"].(string)
		ctx.Set("userID", userID)
		ctx.Set("userRole", claims["role"].(string))

		ctx.Next()
	}
}
