package config

import (
	"database/sql"

	"github.com/minio/minio-go/v7"
)

type Config struct {
	JWTSecret       []byte
	DB              *sql.DB
	MinioClient     *minio.Client
	MinioBucketName string
	FrontEndHost    string
}
