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

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id": id,
	})
}

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
