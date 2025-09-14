package main

import (
	"net/http"

	_ "github.com/alibekkenny/simpengine/cmd/web/docs"
	"github.com/alibekkenny/simpengine/internal/auth"
	"github.com/alibekkenny/simpengine/internal/media"
	romanticevent "github.com/alibekkenny/simpengine/internal/romantic_event"
	simptarget "github.com/alibekkenny/simpengine/internal/simp-target"
	"github.com/alibekkenny/simpengine/internal/user"
	"github.com/justinas/alice"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	userModule := user.NewModule(app.config.DB)
	authModule := auth.NewModule(userModule.Repo)
	simpTargetModule := simptarget.NewModule(app.config.DB)
	romanticEventModule := romanticevent.NewModule(app.config.DB, simpTargetModule.Service)
	mediaModule := media.NewModule(app.config.DB, app.config.MinioClient, app.config.MinioBucketName)

	userModule.RegisterRoutes(mux, app.config)
	authModule.RegisterRoutes(mux, app.config)
	simpTargetModule.RegisterRoutes(mux, app.config)
	romanticEventModule.RegisterRoutes(mux, app.config)
	mediaModule.RegisterRoutes(mux, app.config)

	// chain of middleware
	standardChain := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return standardChain.Then(mux)
}
