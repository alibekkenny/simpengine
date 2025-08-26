package simptarget

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alibekkenny/simpengine/internal/shared/model"
	"github.com/go-playground/validator/v10"
)

type SimpTargetHandler struct {
	service   *SimpTargetService
	validator *validator.Validate
}

func NewSimpTargetHandler(s *SimpTargetService, v *validator.Validate) *SimpTargetHandler {
	return &SimpTargetHandler{service: s, validator: v}
}

func (h *SimpTargetHandler) CreateSimpTarget(w http.ResponseWriter, r *http.Request) {
	var body SimpTargetRequestDTO
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err = h.validator.Struct(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("validation error: %v", err.Error()), http.StatusBadRequest)
		return
	}

	simpTargetID, err := h.service.CreateSimpTarget(r.Context(), body.Name, body.Description)
	if err != nil {
		code, msg := model.ErrorStatus(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", fmt.Sprintf("/simp-target/%s", simpTargetID))
	json.NewEncoder(w).Encode(map[string]any{
		"id":          simpTargetID,
		"name":        body.Name,
		"description": body.Description,
	})
}

func (h *SimpTargetHandler) UpdateSimpTarget(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var body SimpTargetRequestDTO
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err = h.validator.Struct(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("validation error: %v", err), http.StatusBadRequest)
		return
	}

	err = h.service.UpdateSimpTarget(r.Context(), id, body.Name, body.Description)
	if err != nil {
		code, msg := model.ErrorStatus(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SimpTargetHandler) DeleteSimpTarget(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteSimpTarget(r.Context(), id)
	if err != nil {
		code, msg := model.ErrorStatus(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SimpTargetHandler) ViewSimpTarget(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	target, err := h.service.GetSimpTargetByIDAndUser(r.Context(), id)
	if err != nil {
		code, msg := model.ErrorStatus(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(target)
}

func (h *SimpTargetHandler) ViewSimpTargetByUser(w http.ResponseWriter, r *http.Request) {
	targets, err := h.service.GetSimpTargetsByUserID(r.Context())
	if err != nil {
		code, msg := model.ErrorStatus(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(targets)
}
