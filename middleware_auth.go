package main

import (
	"fmt"
	"net/http"

	"github.com/reinhardbuyabo/RSS-Aggregator/internal/auth"
	"github.com/reinhardbuyabo/RSS-Aggregator/internal/database"
)

// this function doesn't match the function signature of a normal http.HandlerFunc
type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// Receiver is a special parameter syntactically placed before the name of the function
func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	// returning a closure ... a new anonymous function ... we'll have access to everything we need from the DB ... through apiCfg *apiConfig
	return func(w http.ResponseWriter, r *http.Request) {

		// 2. LOGIC within auth.go that returns the apiKey of user
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Auth Error: %v", err))
			return
		}

		// 3. QUERYING database by use of apiKey
		// 1. Context
		// 2. apiKey
		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		// 4. calling handler of type authedHandler which handles our request by responding with an authenticated user
		handler(w, r, user)
	}
}
