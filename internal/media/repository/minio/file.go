package minio

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type MinioRepository struct {
	client     *minio.Client
	bucketName string
}

func NewMinioRepository(client *minio.Client, bucketName string) MinioRepository {
	return MinioRepository{client: client, bucketName: bucketName}
}

func (r MinioRepository) Store(ctx context.Context, reader io.Reader, filename string, fileSize int64) (string, error) {
	objectName := uuid.NewString() + "-" + filename

	_, err := r.client.PutObject(ctx, r.bucketName, objectName, reader, fileSize, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (r MinioRepository) Get(ctx context.Context, objectName string) (io.ReadSeekCloser, error) {
	return r.client.GetObject(ctx, r.bucketName, objectName, minio.GetObjectOptions{})
}

// Delete(ctx context.Context, id int64) error
func (r MinioRepository) Delete(ctx context.Context, objectName string) error {
	return r.client.RemoveObject(ctx, r.bucketName, objectName, minio.RemoveObjectOptions{})
}
