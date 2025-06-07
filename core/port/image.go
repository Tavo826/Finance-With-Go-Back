package port

import (
	"context"
	"mime/multipart"
)

type ImageAdapter interface {
	UploadImageFromFile(ctx context.Context, file multipart.File, userId string) (string, error)
	DeleteImage(ctx context.Context, publicId string) error
}
