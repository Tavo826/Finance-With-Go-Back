package main

import (
	"time"
)

type Transaction struct {
	ID               interface{} `json:"_id" bson:"_id,omitempty"`
	Amount           float64     `json:"amount"`
	Type             string      `json:"type"`
	Subject          string      `json:"subject"`
	PersonOrBusiness string      `json:"person_business" bson:"person_business"`
	Description      string      `json:"description"`
	CreatedAtString  string      `json:"created" bson:"created"`
	CreatedAt        time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at,omitempty" bson:"updated_at"`
}
