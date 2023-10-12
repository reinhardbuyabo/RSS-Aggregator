package main

import "net/http"

func handlerErr(w http.ResponseWriter, r *http.Request) {
	// from json.go
	respondWithError(w, 400, "Something went wrong") // Client Error ... link to v1Router := handleFunc("/err", handleErr) ... 2nd Parameter is this function
}
