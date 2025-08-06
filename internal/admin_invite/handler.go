package admininvite

import (
	"encoding/json"
	"net/http"
)

type AdminInviteHandler struct {
	service *AdminInviteService
}

func NewAdminInviteHandler(service *AdminInviteService) *AdminInviteHandler {
	return &AdminInviteHandler{service: service}
}

func (h *AdminInviteHandler) CreateInvite(w http.ResponseWriter, r *http.Request) {
	token, err := h.service.GenerateInvite(r.Context())
	if err != nil {
		http.Error(w, "could not generate token", http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
