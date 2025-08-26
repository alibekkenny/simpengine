package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alibekkenny/simpengine/internal/shared/model"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	service   *AuthService
	validator *validator.Validate
}

func NewAuthHandler(s *AuthService, v *validator.Validate) *AuthHandler {
	return &AuthHandler{service: s, validator: v}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body LoginRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(body); err != nil {
		http.Error(w, fmt.Sprintf("validation error:\n%s", err.Error()), http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), body.Login, body.Password)
	if err != nil {
		code, msg := model.ErrorStatus(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
