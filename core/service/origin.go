package service

import (
	"context"
	"personal-finance/core/domain"
	"personal-finance/core/port"
)

type OriginService struct {
	repo port.OriginRepository
}

func NewOriginService(repo port.OriginRepository) *OriginService {

	return &OriginService{
		repo,
	}
}

func (os *OriginService) GetOriginsByUserId(ctx context.Context, userId string) ([]domain.Origin, error) {

	origins, err := os.repo.GetOriginsByUserId(ctx, userId)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return origins, nil

}

func (os *OriginService) GetOriginById(ctx context.Context, id string) (*domain.Origin, error) {

	origin, err := os.repo.GetOriginById(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		if err.Error() == domain.ErrNoDocuments.Error() {
			return nil, domain.ErrNoDocuments
		}
		return nil, domain.ErrInternal
	}

	return origin, nil
}

func (os *OriginService) CreateOrigin(ctx context.Context, origin *domain.Origin) (*domain.Origin, error) {

	origin, err := os.repo.CreateOrigin(ctx, origin)

	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return origin, nil
}

func (os *OriginService) UpdateOrigin(ctx context.Context, id string, origin *domain.Origin) (*domain.Origin, error) {

	_, err := os.repo.UpdateOrigin(ctx, id, origin)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		if err == domain.ErrDataNotFound {
			return nil, domain.ErrDataNotFound
		}
		return nil, domain.ErrInternal
	}

	return origin, nil
}

func (os *OriginService) DeleteOrigin(ctx context.Context, id string) error {

	return os.repo.DeleteOrigin(ctx, id)
}
