package dto

import (
	"personal-finance/core/domain"
	"time"
)

type User struct {
	ID           string    `json:"_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	Role         string    `json:"role"`
	ProfileImage string    `json:"profile_image,omitempty"`
	PublicId     string    `json:"public_id,omitempty"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRequest struct {
	ID string `form:"id" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func NewUserResponse(user *domain.User) User {

	return User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Role:         user.Role,
		ProfileImage: user.ProfileImage,
		PublicId:     user.PublicIdImage,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
