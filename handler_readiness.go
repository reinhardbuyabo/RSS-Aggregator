package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct{}{}) // empty struct which will marshall into an empty json ... that's the payload.
	// hooking up a handler ... creating v1Router

} // Specific function signature for http handler as per Go's standard library expects
