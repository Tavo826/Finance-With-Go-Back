package dto

import (
	"personal-finance/core/domain"
	"time"
)

type OriginByUserId struct {
	UserId string `form:"user_id" binding:"required"`
}

type OriginRequest struct {
	UserId    string    `json:"user_id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	Total     float64   `json:"total" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type OriginResponse struct {
	ID        string    `json:"_id" binding:"required"`
	UserId    string    `json:"user_id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	Total     float64   `json:"total" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func NewOriginResponse(origin *domain.Origin) OriginResponse {

	return OriginResponse{
		ID:        origin.ID,
		UserId:    origin.UserId,
		Name:      origin.Name,
		Total:     origin.Total,
		CreatedAt: origin.CreatedAt,
		UpdatedAt: origin.UpdatedAt,
	}
}
