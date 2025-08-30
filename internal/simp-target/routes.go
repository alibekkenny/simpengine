package simptarget

import (
	"database/sql"
	"net/http"

	"github.com/alibekkenny/simpengine/cmd/config"
	"github.com/alibekkenny/simpengine/internal/auth"
	"github.com/go-playground/validator/v10"
)

type Module struct {
	Service *SimpTargetService
}

func NewModule(db *sql.DB) *Module {
	repo := NewPosgresRepository(db)
	service := NewSimpTargetService(repo)
	return &Module{Service: service}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux, cfg *config.Config) {
	validator := validator.New()
	handler := NewSimpTargetHandler(m.Service, validator)

	mux.Handle("POST /simp-target", auth.AuthMiddleware(http.HandlerFunc(handler.CreateSimpTarget)))
	mux.Handle("PUT /simp-target/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.UpdateSimpTarget)))
	mux.Handle("DELETE /simp-target/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.DeleteSimpTarget)))
	mux.Handle("GET /simp-target/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.ViewSimpTarget)))
	mux.Handle("GET /simp-target", auth.AuthMiddleware(http.HandlerFunc(handler.ViewSimpTargetByUser)))
}
