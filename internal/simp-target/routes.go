package simptarget

import (
	"net/http"

	"github.com/alibekkenny/simpengine/cmd/config"
	"github.com/alibekkenny/simpengine/internal/auth"
	"github.com/go-playground/validator/v10"
)

func RegisterRoutes(mux *http.ServeMux, cfg *config.Config) {
	repo := NewPosgresRepository(cfg.DB)
	service := NewSimpTargetService(repo)
	validator := validator.New()
	handler := NewSimpTargetHandler(service, validator)

	mux.Handle("GET /simp-target/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.ViewSimpTarget)))
}
