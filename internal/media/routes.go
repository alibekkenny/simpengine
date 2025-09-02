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

func NewModule(db *sql.DB, minioEndpoint string, bucketName string) (*Module, error) {
	client, err := minio.New(minioEndpoint, &minio.Options{})
	if err != nil {
		return nil, err
	}

	fileRepo := minio_repo.NewMinioRepository(client, bucketName)
	mediaRepo := postgres.NewMediaRepository(db)
	service := NewMediaService(mediaRepo, fileRepo)

	return &Module{Service: service}, nil
}

func (m *Module) RegisterRoutes(mux *http.ServeMux, cfg *config.Config) {
	handler := NewMediaHandler(m.Service)
	mux.Handle("POST /media", auth.AuthMiddleware(http.HandlerFunc(handler.UploadFile)))
	mux.Handle("GET /media/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.DownloadFile)))
	// mux.Handler("DELETE /media/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.DeleteFile)))
}
