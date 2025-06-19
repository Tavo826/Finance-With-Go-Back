package domain

import "time"

type Origin struct {
	ID        string    `json:"_id" bson:"_id,omitempty"`
	UserId    string    `json:"user_id" bson:"user_id" validate:"required"`
	Name      string    `json:"name" bson:"name" validate:"required"`
	Total     float64   `json:"total" bson:"total" validate:"required"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

type OriginRequest struct {
	UserId    string    `json:"user_id" bson:"user_id" validate:"required"`
	Name      string    `json:"name" bson:"name" validate:"required"`
	Total     float64   `json:"total" bson:"total" validate:"required"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}
