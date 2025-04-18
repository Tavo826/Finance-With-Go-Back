package main

import (
	"time"
)

type Transaction struct {
	ID               any       `json:"_id" bson:"_id,omitempty"`
	Amount           float64   `json:"amount" validate:"required"`
	Type             string    `json:"type" validate:"required"`
	Subject          string    `json:"subject" validate:"required"`
	PersonOrBusiness string    `json:"person_business" bson:"person_business" validate:"required"`
	Description      string    `json:"description" validate:"required"`
	CreatedAtString  string    `json:"created" bson:"created" validate:"required"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

type Error struct {
	Key   string `json:"key"`
	Error string `json:"error"`
}

type ErrorResponse struct {
	message      string  `json:"message"`
	errorMessage []Error `json:"error"`
}
