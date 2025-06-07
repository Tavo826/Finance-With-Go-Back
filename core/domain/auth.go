package domain

import "time"

type User struct {
	ID            string    `json:"_id" bson:"_id,omitempty"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	Role          string    `json:"role"`
	ProfileImage  string    `json:"profile_image"`
	PublicIdImage string    `json:"public_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Image struct {
	SecureUrl string `json:"secure_url"`
	PublicId  string `json:"public_id"`
}
