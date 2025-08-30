package auth

import (
	"net/http"

	"github.com/alibekkenny/simpengine/cmd/config"
	"github.com/alibekkenny/simpengine/internal/user"
	"github.com/go-playground/validator/v10"
)

type Module struct {
	Service *AuthService
}

func NewModule(repo user.UserRepository) *Module {
	service := NewAuthService(repo)
	return &Module{Service: service}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux, cfg *config.Config) {
	InitJWT(cfg.JWTSecret)

	validator := validator.New()
	handler := NewAuthHandler(m.Service, validator)

	mux.HandleFunc("POST /user/login", handler.Login)
}
