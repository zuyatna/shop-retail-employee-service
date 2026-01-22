package usecase

import (
	"context"
	"io"
)

type StorageRepository interface {
	UploadFile(ctx context.Context, fileName string, contentType string, content io.Reader, size int64) (string, error)
}
