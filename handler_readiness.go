package main

import "net/http"

// Just responds if server is alive and running
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	// 1. ResponseWriter
	// 2. Pointer to http request
	// from json.go
	respondWithJSON(w, 200, struct{}{}) // passing writer, responding with 200, and some reponse payload, ... , empty struct marshalls to JSON
}
