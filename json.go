package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) { // format message into a consistent json everytime
	if code > 499 { // logging a message with server errors ... this means we have a bug on our end
		log.Println("Responding with 5XX error: ", msg)
	}
	type errResponse struct { // So below we have an error field with an error key.
		Error string `json:"error"` // here we are saying the key that this should marshall to is error ... we're specifying how we want that json.marshall and json.umarshall to convert the struct to a marshall object
	}

	// Looks like something like this
	// {
	// 	"error": "something went wrong"
	// }

	respondWithJSON(w, code, errResponse{
		Error: msg,
	}) // responding with a specific structure
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// 1. http handlers in go use this write
	// 2. status code to respond with
	// 3. Sth that we can marshall to a json structure

	data, err := json.Marshal(payload) // marshall whatever it's given(payload) into a json string and it will return it as bytes. So that we can write it in a binary format directly to the http request
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500) // internal server error
		return
	}
	w.Header().Add("Content-Type", "application/json") // adding a response to the response header: we're responding with content type json which is the standard value for json responses
	w.WriteHeader(code)                                // ok status code .... Everything went well
	w.Write(data)                                      // Writing the response body
	// Now we have a way to respond with some JSON data ... now we have to create a handler
}
