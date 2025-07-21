package main

import (
	"net/http"
	"log"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", neuter(fileServer)))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/simp-target/view", simpTargetView)
	mux.HandleFunc("/simp-target/create", simpTargetCreate)

	
	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

func neuter(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if strings.HasSuffix(r.URL.Path, "/") {
            http.NotFound(w, r)
            return
        }

        next.ServeHTTP(w, r)
    })
}