package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service   *UserService
	validator *validator.Validate
}

func NewUserHandler(s *UserService, v *validator.Validate) *UserHandler {
	return &UserHandler{service: s, validator: v}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body RegisterRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(body); err != nil {
		http.Error(w, fmt.Sprintf("Validation error:\n%v", err), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	id, err := h.service.Register(ctx, body.Login, body.Email, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/user/%s", id))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    id,
		"login": body.Login,
		"email": body.Email,
	})
}
