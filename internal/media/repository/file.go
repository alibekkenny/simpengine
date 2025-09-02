package repository

import (
	"context"
	"io"
)

type FileRepository interface {
	Store(ctx context.Context, reader io.Reader, filename string, fileSize int64) (string, error)
	Get(ctx context.Context, objectName string) (io.ReadSeekCloser, error)
	Delete(ctx context.Context, objectName string) error
}
