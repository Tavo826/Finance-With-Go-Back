package port

import (
	"context"
	"mime/multipart"
	"personal-finance/core/domain"
)

type ImageAdapter interface {
	UploadImageFromFile(ctx context.Context, file multipart.File, userId string) (*domain.Image, error)
	DeleteImage(ctx context.Context, publicId string) error
}
