package cloudinary

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/ingwrok/hotelBooking/internal/common/errs"
	"github.com/ingwrok/hotelBooking/internal/common/logger"
	"go.uber.org/zap"
)

type CloudinaryImageUploader struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryImageUploader(cloudName, apiKey, apiSecret string) (*CloudinaryImageUploader, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}
	return &CloudinaryImageUploader{cld: cld}, nil
}

func (u *CloudinaryImageUploader) UploadImage(ctx context.Context, file io.Reader, filename string) (string, error) {
	logger.Info("Uploading image to Cloudinary", zap.String("filename", filename))

	resp, err := u.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:   "hotel_booking_room_types",
		PublicID: filename,
	})

	if err != nil {
		logger.ErrorErr(err, "Cloudinary upload failed")
		return "", errs.NewUnexpectedError("image upload failed")
	}

	logger.Info("Image uploaded successfully", zap.String("url", resp.SecureURL))
	return resp.SecureURL, nil
}
