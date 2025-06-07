package cloud

import (
	"context"
	"log/slog"
	"personal-finance/adapter/config"

	"github.com/cloudinary/cloudinary-go/v2"
)

func New(ctx context.Context, config *config.ImageCloud) (*cloudinary.Cloudinary, error) {

	cld, err := cloudinary.New()
	if err != nil {
		slog.Error("error connecting to image service", "error", err)
		return nil, err
	}

	return cld, nil
}
