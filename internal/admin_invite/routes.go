package admininvite

import (
	"net/http"

	"github.com/alibekkenny/simpengine/cmd/config"
)

func RegisterRoutes(mux *http.ServeMux, config *config.Config) {
	repo := NewPosgresRepository(config.DB)
	service := NewAdminInviteService(repo)
	handler := NewAdminInviteHandler(service)

	mux.HandleFunc("POST /admin/invite", handler.CreateInvite)
}
