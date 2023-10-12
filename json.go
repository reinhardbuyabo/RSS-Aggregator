package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error: ", msg)
	}

	type errResponse struct {
		Error string `json:"error"` // we take ... to convert this struct to a json  that looks like {"error": "something went wrong"}
	}

	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// payload: sth we can marshall to json struct ...

	data, err := json.Marshal(payload) // return as bytes ... binary format

	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500) // internal server error
		return
	}

	// Adding Header
	w.Header().Add("Content-Type", "application/json") // adding response header saying we're responding with JSON
	w.WriteHeader(code)
	w.Write(data) // response body
}
