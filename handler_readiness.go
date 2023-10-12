package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct{}{}) // passing writer, responding with 200, and some payload, ... , marshalls to JSON
}
