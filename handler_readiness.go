package main

import "net/http"

// Just responds if server is alive and running
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	// from json.go
	respondWithJSON(w, 200, struct{}{}) // passing writer, responding with 200, and some payload, ... , marshalls to JSON
}
