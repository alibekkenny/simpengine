package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)

	mux.HandleFunc("/simp-target/view", app.simpTargetView)
	mux.HandleFunc("/simp-target/create", app.simpTargetCreate)

	mux.HandleFunc("/user/view", app.userView)
	mux.HandleFunc("/user/create", app.userCreate)

	return mux
}
