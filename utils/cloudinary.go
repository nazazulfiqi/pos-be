package utils

import (
	"context"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func Bool(v bool) *bool {
	return &v
}

func UploadToCloudinary(file multipart.File, fileName string) (string, string, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return "", "", err
	}

	ctx := context.Background()

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:  "products/" + fileName,
		Overwrite: Bool(true), // ✅ sekarang pointer ke bool
	})
	if err != nil {
		return "", "", err
	}

	return uploadResult.SecureURL, uploadResult.PublicID, nil
}

func DeleteFromCloudinary(publicID string) error {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
	return err
}
