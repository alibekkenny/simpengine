package user

import (
	"net/http"

	"github.com/alibekkenny/simpengine/cmd/config"
	"github.com/go-playground/validator/v10"
)

func RegisterRoutes(mux *http.ServeMux, config *config.Config) {
	repo := NewPosgresRepository(config.DB)
	service := NewUserService(repo)
	validator := validator.New()

	handler := NewUserHandler(service, validator)

	mux.HandleFunc("POST /user/register", handler.Register)
}
