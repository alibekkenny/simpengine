package media

import (
	"database/sql"
	"net/http"

	"github.com/alibekkenny/simpengine/cmd/config"
	"github.com/alibekkenny/simpengine/internal/auth"
	minio_repo "github.com/alibekkenny/simpengine/internal/media/repository/minio"
	"github.com/alibekkenny/simpengine/internal/media/repository/postgres"
	"github.com/minio/minio-go/v7"
)

type Module struct {
	Service *MediaService
}

func NewModule(db *sql.DB, minioClient *minio.Client, bucketName string) *Module {
	fileRepo := minio_repo.NewMinioRepository(minioClient, bucketName)
	mediaRepo := postgres.NewMediaRepository(db)
	service := NewMediaService(mediaRepo, fileRepo)

	return &Module{Service: service}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux, cfg *config.Config) {
	handler := NewMediaHandler(m.Service)
	mux.Handle("POST /media", auth.AuthMiddleware(http.HandlerFunc(handler.UploadFile)))
	mux.Handle("GET /media/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.DownloadFile)))
	mux.Handle("DELETE /media/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.DeleteFile)))
}
