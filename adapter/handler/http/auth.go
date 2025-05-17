package http

import (
	"log"
	"personal-finance/adapter/config"
	"personal-finance/adapter/handler/http/dto"
	"personal-finance/core/domain"
	"personal-finance/core/port"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	service  port.AuthService
	validate *validator.Validate
	config   *config.Token
}

func NewAuthHandler(service port.AuthService, validate *validator.Validate, config *config.Token) *AuthHandler {

	return &AuthHandler{
		service,
		validate,
		config,
	}
}

func (ah *AuthHandler) Register(ctx *gin.Context) {

	var req dto.RegisterRequest
	jwtSecret := []byte(ah.config.JwtSecret)

	if err := ctx.Bind(&req); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	if err := ah.validate.Struct(req); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	exists, err := ah.service.VerifyUserEmail(ctx, req.Email)
	if err != nil {
		log.Println("Error verify email exists")
		dto.HandleError(ctx, err)
		return
	}

	if exists {
		dto.HandleError(ctx, domain.ErrUserAlreadyExists)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashed password")
		dto.HandleError(ctx, domain.ErrInvalidToken)
		return
	}

	user := domain.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = ah.service.CreateUser(ctx, &user)
	if err != nil {
		dto.HandleError(ctx, err)
		return
	}

	token, err := generateToken(&user, jwtSecret)
	if err != nil {
		dto.HandleError(ctx, domain.ErrInternal)
		return
	}

	tokenResponse := dto.TokenResponse{
		Token: token,
		User:  dto.NewUserResponse(&user),
	}

	dto.HandleSuccess(ctx, tokenResponse)
}

func (ah *AuthHandler) Login(ctx *gin.Context) {

	var req dto.LoginRequest
	jwtSecret := []byte(ah.config.JwtSecret)

	if err := ctx.Bind(&req); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	if err := ah.validate.Struct(req); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	user, err := ah.service.GetUserByEmail(ctx, req.Email)
	if err != nil {
		dto.HandleError(ctx, domain.ErrDataNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		dto.HandleError(ctx, domain.ErrUnauthorized)
		return
	}

	token, err := generateToken(user, jwtSecret)
	if err != nil {
		dto.HandleError(ctx, domain.ErrInternal)
		return
	}

	tokenResponse := dto.TokenResponse{
		Token: token,
		User:  dto.NewUserResponse(user),
	}

	dto.HandleSuccess(ctx, tokenResponse)
}

func generateToken(user *domain.User, jwtSecret []byte) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
