package adapter

import (
	"context"
	"encoding/json"
	"log"
	"mime/multipart"
	"personal-finance/adapter/storage/cloud/data"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type ImageAdapter struct {
	adapter *cloudinary.Cloudinary
}

func NewImageAdapter(adapter *cloudinary.Cloudinary) *ImageAdapter {
	return &ImageAdapter{
		adapter: adapter,
	}
}

func (ia *ImageAdapter) UploadImageFromFile(ctx context.Context, file multipart.File, userId string) (string, error) {

	var imageResponse data.ImageUploadResponse

	file.Seek(0, 0)

	overwrite := true

	uploadParams := uploader.UploadParams{
		Folder:         "profile_images",
		PublicID:       "user_" + userId,
		Overwrite:      &overwrite,
		ResourceType:   "image",
		Transformation: "c_fill,w_300,h_300,q_auto,f_auto",
	}

	response, err := ia.adapter.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		log.Println("Response: ", response)
		log.Println("Error: ", err)
		return "", err
	}

	resp, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(resp, &imageResponse); err != nil {
		return "", err
	}

	return imageResponse.SecureUrl, nil
}

func (ia *ImageAdapter) DeleteImage(ctx context.Context, publicId string) error {
	_, err := ia.adapter.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicId,
	})

	return err
}
