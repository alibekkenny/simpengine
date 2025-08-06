package main

import (
	"net/http"

	adminInvite "github.com/alibekkenny/simpengine/internal/admin_invite"
	"github.com/alibekkenny/simpengine/internal/auth"
	simptarget "github.com/alibekkenny/simpengine/internal/simp-target"
	"github.com/alibekkenny/simpengine/internal/user"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	user.RegisterRoutes(mux, app.config)
	adminInvite.RegisterRoutes(mux, app.config)
	auth.RegisterRoutes(mux, app.config)
	simptarget.RegisterRoutes(mux, app.config)

	// chain of middleware
	standardChain := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standardChain.Then(mux)
}
