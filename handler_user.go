package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/reinhardbuyabo/RSS-Aggregator/internal/database"
)

// CREATE USER HANDLER ... action in response to a route

// We want to add a function ... function signature remains the same, now we have an additional data
func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Name string `json:"name"`
	}

	// parsing request body to struct ... returns a decoder
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	// we want to decode into an instance of the parameter struct
	err := decoder.Decode(&params) // decode into an instance of the parameter struct
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// Using database to create user ... method that sqlc generated for us ... created parameters as a struct
	// Parameters
	// 1. Context
	// 2. Create User Parameters

	// CreateUser() Returns
	// New User or Error
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(), // uuid - really long id represented as text
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name, // whatever was passed in the Name params ... that is
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	// 1. ResponseWriter
	// 2. Status Code
	// 3. Database User
	// from json.go

	respondWithJSON(w, 201, databaseUserToUser(user)) // passing writer, responding with 200, and some reponse payload, ... , empty struct marshalls to JSON
}

// This is an authenticated endpoint
func (apiCfg *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request, user database.User) { // handleGetUser() doesn't match the function signature of handleFunc, but matches the function signature of authedHandler hence we can't pass it as the 2nd parameter to v1Router.Get(), but we can pass it to middlewareAuth() which returns a http.HandlerFunc
	// apiKey, err := auth.GetAPIKey(r.Header)
	// if err != nil {
	// 	respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
	// }

	// // 1. Request's Context - Context Package in the go standard library, it gives you a way to track some thing that is happening in multiple go routines ... you can cancel Context, effectively killing a http Request
	// // 2. apiKey
	// user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
	// if err != nil {
	// 	respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
	// 	return
	// }

	respondWithJSON(w, 200, databaseUserToUser(user))
}
