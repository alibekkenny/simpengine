package media

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/alibekkenny/simpengine/internal/shared/model"
)

type MediaHandler struct {
	service *MediaService
}

func NewMediaHandler(service *MediaService) *MediaHandler {
	return &MediaHandler{service: service}
}

// UploadFile godoc
// @Summary      Upload a media file
// @Description  Uploads a file to the media storage and returns its ID.
// @Tags         media
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "File to upload"
// @Success      201   {object}  UploadMediaResponseDTO
// @Failure      400   {object}  model.ErrorResponse
// @Failure      401   {object}  model.ErrorResponse
// @Failure      500   {object}  model.ErrorResponse
// @Security     BearerAuth
// @Router       /media [post]
func (h *MediaHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid form: %v", err), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid form: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	id, err := h.service.UploadFile(r.Context(), file, header.Filename, header.Size)
	if err != nil {
		model.WriteErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(UploadMediaResponseDTO{
		ID: id,
	})
}

// DownloadFile  downloads a media file by ID.
// @Summary      Download media file
// @Description  Returns the file as binary stream (with correct Content-Type and Content-Disposition headers).
// @Tags         media
// @Param        id   path      int  true  "Media ID"
// @Produce      application/octet-stream
// @Success      200  {file}    file
// @Failure      400  {object}  model.ErrorResponse  "Invalid ID"
// @Failure      404  {object}  model.ErrorResponse  "File not found"
// @Failure      500  {object}  model.ErrorResponse  "Internal server error"
// @Security     BearerAuth
// @Router       /media/{id} [get]
func (h *MediaHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	media, object, err := h.service.DownloadFile(r.Context(), id)
	if err != nil {
		model.WriteErrorResponse(w, err)
		return
	}
	defer object.Close()

	w.Header().Set("Content-Type", media.MimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", media.OriginalName))
	if _, err := io.Copy(w, object); err != nil {
		http.Error(w, "Error streaming file", http.StatusInternalServerError)
		return
	}
}

// DeleteFile  	 deletes a media file by ID.
// @Summary      Deletes media file
// @Description  Deletes the file by its id. Only media owner can delete it.
// @Tags         media
// @Param        id   path      int  true  "Media ID"
// @Success      204  "No content"
// @Failure      400  {object}  model.ErrorResponse  "Invalid ID"
// @Failure      404  {object}  model.ErrorResponse  "File not found"
// @Failure      500  {object}  model.ErrorResponse  "Internal server error"
// @Security     BearerAuth
// @Router       /media/{id} [delete]
func (h *MediaHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteFile(r.Context(), id); err != nil {
		model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
