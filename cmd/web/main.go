package main

import (
	"net/http"
	"log"
)

func main() {
	mux := http.NewServeMux()
	
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/simp-target/view", simpTargetView)
	mux.HandleFunc("/simp-target/create", simpTargetCreate)

	
	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}