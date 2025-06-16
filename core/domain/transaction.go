package domain

import (
	"time"
)

type Transaction struct {
	ID               string    `json:"_id" bson:"_id,omitempty"`
	UserId           string    `json:"user_id" bson:"user_id" validate:"required"`
	OriginId         *string   `json:"origin_id,omitempty" bson:"origin_id,omitempty"`
	Amount           float64   `json:"amount" validate:"required"`
	Type             string    `json:"type" validate:"required"`
	Subject          string    `json:"subject" validate:"required"`
	PersonOrBusiness string    `json:"person_business" bson:"person_business" validate:"required"`
	Description      string    `json:"description" validate:"required"`
	CreatedAtString  string    `json:"created" bson:"created" validate:"required"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at,omitempty" bson:"updated_at"`
	Origin           *Origin   `json:"origin,imitempty" bson:"origin,omitempty"`
}
