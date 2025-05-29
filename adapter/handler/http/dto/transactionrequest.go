package dto

import "time"

type TransactionByUserRequest struct {
	UserId string `form:"user_id" binding:"required"`
	Page   uint64 `form:"page" binding:"required"`
	Limit  uint64 `form:"limit" binding:"required"`
}

type DateFilterRequest struct {
	*TransactionByUserRequest
	Month int `form:"month" binding:"min=0,max=12"`
	Year  int `form:"year" binding:"required"`
}

type SubjectFilterRequest struct {
	*TransactionByUserRequest
	Subject          string `form:"subject" binding:"required"`
	PersonOrBusiness string `form:"person_business"`
}

type IdRequest struct {
	ID string `uri:"id" binding:"required"`
}

type TransactionRequest struct {
	Amount           float64   `json:"amount" validate:"required" binding:"gte=0"`
	UserId           string    `json:"user_id" bson:"user_id" validate:"required"`
	Type             string    `json:"type" validate:"required"`
	Subject          string    `json:"subject" validate:"required"`
	PersonOrBusiness string    `json:"person_business" bson:"person_business" validate:"required"`
	Description      string    `json:"description" validate:"required"`
	CreatedAtString  string    `json:"created" bson:"created" validate:"required"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
}
