package user

import (
	"database/sql"
	"net/http"

	"github.com/alibekkenny/simpengine/cmd/config"
	"github.com/go-playground/validator/v10"
)

type Module struct {
	Service *UserService
	Repo    UserRepository
}

func NewModule(db *sql.DB) *Module {
	repo := NewPosgresRepository(db)
	service := NewUserService(repo)

	return &Module{Service: service, Repo: repo}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux, cfg *config.Config) {
	validator := validator.New()

	handler := NewUserHandler(m.Service, validator)

	mux.HandleFunc("POST /user/register", handler.Register)
}
