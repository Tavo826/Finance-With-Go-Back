package http

import (
	"log"
	"net/http"
	"personal-finance/adapter/config"
	"personal-finance/adapter/handler/http/dto"
	"personal-finance/core/domain"
	"personal-finance/core/port"
	"strings"
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

func (ah *AuthHandler) GetUserById(ctx *gin.Context) {
	var request dto.UserRequest

	log.Println("Handler getting user")

	if err := ctx.Bind(&request); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	user, err := ah.service.GetUserById(ctx, request.ID)
	if err != nil {
		dto.HandleError(ctx, err)
		return
	}

	response := dto.NewUserResponse(user)

	dto.HandleSuccess(ctx, response)
}

func (ah *AuthHandler) UpdateUser(ctx *gin.Context) {

	log.Println("Update request: ", ctx.Request)

	const MaxImageSize = 5 << 20
	id := ctx.Param("id")

	if err := ctx.Request.ParseMultipartForm(MaxImageSize); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	file, header, err := ctx.Request.FormFile("profile_image")
	if err != nil && err != http.ErrMissingFile {
		dto.HandleError(ctx, domain.ErrGettingFile)
		return
	}

	username := ctx.PostForm("username")
	email := ctx.PostForm("email")

	user := domain.User{
		Username:  username,
		Email:     email,
		UpdatedAt: time.Now(),
	}

	var uploadedImage *domain.Image

	if file != nil {
		defer file.Close()

		if header != nil && header.Size > MaxImageSize {
			dto.HandleError(ctx, domain.ErrFileSize)
			return
		}

		contentType := header.Header.Get("Content-Type")
		if !isValidImageType(contentType) {
			dto.HandleError(ctx, domain.ErrFileType)
			return
		}

		uploadedImage, err = ah.service.UpdateUserProfileImage(ctx, file, id)
		if err != nil {
			dto.HandleError(ctx, err)
			return
		}

		user.ProfileImage = uploadedImage.SecureUrl
		user.PublicIdImage = uploadedImage.PublicId
	}

	actualUser, err := ah.service.GetUserById(ctx, id)
	if err != nil {
		dto.HandleError(ctx, err)
		return
	}

	user.Password = actualUser.Password
	user.Role = actualUser.Role
	user.CreatedAt = actualUser.CreatedAt

	_, err = ah.service.UpdateUser(ctx, id, &user)
	if err != nil {
		dto.HandleError(ctx, err)
		return
	}

	response := dto.NewUserResponse(&user)

	dto.HandleSuccess(ctx, response)
}

func (ah *AuthHandler) DeleteUser(ctx *gin.Context) {
	var request dto.IdRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		dto.ValidationError(ctx, err)
		return
	}

	err := ah.service.DeleteUser(ctx, request.ID)
	if err != nil {
		dto.HandleError(ctx, err)
		return
	}

	err = ah.service.DeleteTransactionsByUserId(ctx, request.ID)
	if err != nil {
		dto.HandleError(ctx, err)
		return
	}

	dto.HandleSuccess(ctx, nil)
}

func generateToken(user *domain.User, jwtSecret []byte) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
	}

	for _, validType := range validTypes {
		if strings.EqualFold(contentType, validType) {
			return true
		}
	}
	return false
}
