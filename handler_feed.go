package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/reinhardbuyabo/RSS-Aggregator/internal/database"
)

// CREATE FEED HANDLER

// authenticated Endpoint
func (apiCfg *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	// Struct to Reading from Request Body
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	// We need a new decoder that reads from http.Request `r`
	decoder := json.NewDecoder(r.Body)

	// instantiating an object of the parameters struct above
	params := parameters{}

	// Decoding
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %V", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't Create Feed: %v", err))
		return
	}

	// 1. ResponseWriter
	// 2. Status Code
	// 3. Database Feed
	respondWithJSON(w, 201, databaseFeedtoFeed(feed))
}
