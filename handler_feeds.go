package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ab-elhaddad/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerFeedsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	p := parameters{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&p)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding JSON: %v", err))
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		Name:   p.Name,
		Url:    p.Url,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating feed: %v", err))
	}

	respondWithJSON(w, 201, databaseFeedToFeed(feed))
}

func (apiCfg *apiConfig) handlerAllFeedsGet(w http.ResponseWriter, r *http.Request, _user database.User) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting feeds: %v", err))
	}

	feedsToReturn := []Feed{}
	for _, feed := range feeds {
		feedsToReturn = append(feedsToReturn, databaseFeedToFeed(feed))
	}
	respondWithJSON(w, 200, feedsToReturn)
}

func (apiCfg *apiConfig) handlerGetFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	feedId, err := uuid.Parse(chi.URLParam(r, "feedId"))
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing feedId: %v", err))
		return
	}

	feed, err := apiCfg.DB.GetFeedByID(r.Context(), feedId)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting feed: %v", err))
		return
	}

	respondWithJSON(w, 200, databaseFeedToFeed(feed))
}

func (apiCfg *apiConfig) handlerGetFeedsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := apiCfg.DB.GetFeedsByUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting feeds: %v", err))
		return
	}

	feedsToReturn := []Feed{}
	for _, feed := range feeds {
		feedsToReturn = append(feedsToReturn, databaseFeedToFeed(feed))
	}
	respondWithJSON(w, 200, feedsToReturn)
}
