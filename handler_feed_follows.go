package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/reinhardbuyabo/RSS-Aggregator/internal/database"
)

// Authenticated Endpoint/Handler
func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	// 1. Create struct that contains fields with required data from the user
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	// 2. Decode user input given from r.Body
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// 3. Use the decoded parameters to populate database
	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't Create Feed Follow: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedFollowtoFeedFollow(feedFollow))
}

// Authenticated Endpoint
func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	// // 1. Parameters Struct
	// type parameters struct {
	// 	UserID uuid.UUID `json:"user_id"`
	// }
	// // 2. Decode
	// decoder := json.NewDecoder(r.Body)
	// params := parameters{}
	// err := decoder.Decode(&params)
	// if err != nil {
	// 	respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
	// 	return
	// }
	// 3. Response
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't Get Feed Follows, %v", err))
		return
	}

	respondWithJSON(w, 200, databaseFeedFollowsToFeedFollows(feedFollows))
}

// Authenticated Endpoint
func (apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	// This action/route doesn't need payload

	// No need of a body in the payload

	// having the url parameter as the feedFollowID ... we want the feedFollowID dynamically passed ... how
	// 1, Request `r`
	// 2. key => 1st parameter
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr) // returns
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse feed follow id: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete feed follow: %v", err))
		return
	}

	// responding with an empty JSON object
	respondWithJSON(w, 200, struct{}{})
}
