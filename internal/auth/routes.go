package auth

import (
	"net/http"

	"github.com/alibekkenny/simpengine/cmd/config"
	"github.com/alibekkenny/simpengine/internal/user"
	"github.com/go-playground/validator/v10"
)

func RegisterRoutes(mux *http.ServeMux, cfg *config.Config) {
	InitJWT(cfg.JWTSecret)

	repo := user.NewPosgresRepository(cfg.DB)
	service := NewAuthService(repo)
	validator := validator.New()
	handler := NewAuthHandler(service, validator)

	mux.HandleFunc("POST /user/login", handler.Login)
}
