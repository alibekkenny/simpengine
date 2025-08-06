package simptarget

import (
	"net/http"

	"github.com/alibekkenny/simpengine/cmd/config"
	"github.com/alibekkenny/simpengine/internal/auth"
)

func RegisterRoutes(mux *http.ServeMux, cfg *config.Config) {
	handler := NewSimpTargetHandler()

	mux.Handle("GET /simp-target", auth.AuthMiddleware(http.HandlerFunc(handler.ViewSimpTarget)))
}
