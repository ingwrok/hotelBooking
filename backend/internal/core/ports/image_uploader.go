package ports

import (
	"context"
	"io"
)

type ImageUploader interface {
	UploadImage(ctx context.Context, file io.Reader, filename string) (string, error)
}
