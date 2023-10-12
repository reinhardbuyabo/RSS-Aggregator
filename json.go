package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// payload: sth we can marshall to json struct ...

	data, err := json.Marshal(payload) // return as bytes ... binary format

	if err != nil {
		log.Println("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500) // internal server error
		return
	}

	// Adding Header
	w.Header().Add("Content-Type", "application/json") // adding response header saying we're responding with JSON
	w.WriteHeader(200)
	w.Write(data) // response body
}
