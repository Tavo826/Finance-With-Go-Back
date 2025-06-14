package port

import (
	"context"
	"personal-finance/core/domain"
)

type OriginRepository interface {
	GetOriginsByUserId(ctx context.Context, userId string) ([]domain.Origin, error)
	GetOriginById(ctx context.Context, id string) (*domain.Origin, error)
	CreateOrigin(ctx context.Context, origin *domain.Origin) (*domain.Origin, error)
	UpdateOrigin(ctx context.Context, id string, updatedOrigin *domain.Origin) (*domain.Origin, error)
	DeleteOrigin(ctx context.Context, id string) error
}

type OriginService interface {
	GetOriginsByUserId(ctx context.Context, userId string) ([]domain.Origin, error)
	GetOriginById(ctx context.Context, id string) (*domain.Origin, error)
	CreateOrigin(ctx context.Context, origin *domain.Origin) (*domain.Origin, error)
	UpdateOrigin(ctx context.Context, id string, origin *domain.Origin) (*domain.Origin, error)
	DeleteOrigin(ctx context.Context, id string) error
}
