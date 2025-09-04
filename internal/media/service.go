package media

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/alibekkenny/simpengine/internal/auth"
	media_model "github.com/alibekkenny/simpengine/internal/media/model"
	"github.com/alibekkenny/simpengine/internal/media/repository"
	"github.com/alibekkenny/simpengine/internal/shared/model"
)

type MediaService struct {
	mediaRepo repository.MediaRepository
	fileRepo  repository.FileRepository
}

func NewMediaService(mediaRepo repository.MediaRepository, fileRepo repository.FileRepository) *MediaService {
	return &MediaService{mediaRepo: mediaRepo, fileRepo: fileRepo}
}

func (s *MediaService) UploadFile(ctx context.Context, file multipart.File, filename string, fileSize int64) (int64, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return 0, model.ErrInvalidCredentials
	}

	objectName, err := s.fileRepo.Store(ctx, file, filename, fileSize)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	mimeType, err := GetMimeType(file)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	id, err := s.mediaRepo.Create(ctx, objectName, filename, mimeType, fileSize, userID)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return 0, fmt.Errorf("%w: user not found", model.ErrNoRecord)
		}
		return 0, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return id, nil
}

func (s *MediaService) DownloadFile(ctx context.Context, id int64) (*media_model.Media, io.ReadSeekCloser, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return nil, nil, model.ErrInvalidCredentials
	}

	media, err := s.mediaRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return nil, nil, fmt.Errorf("%w: file not found", model.ErrNoRecord)
		}
		return nil, nil, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	object, err := s.fileRepo.Get(ctx, media.ObjectName)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return media, object, nil
}

func (s *MediaService) DeleteFile(ctx context.Context, id int64) error {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return model.ErrInvalidCredentials
	}

	media, err := s.mediaRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: file not found", model.ErrNoRecord)
		}

		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	if err := s.mediaRepo.Delete(ctx, id, userID); err != nil {
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	if err := s.fileRepo.Delete(ctx, media.ObjectName); err != nil {
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return nil
}
