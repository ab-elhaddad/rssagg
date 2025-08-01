package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ab-elhaddad/rssagg/internal/database"
)

func (apiCfg *apiConfig) handlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	p := parameters{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&p)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Name:     p.Name,
		Email:    p.Email,
		Password: p.Password,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	user.Password = ""
	respondWithJSON(w, 201, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	p := parameters{}

	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()
	err := decoder.Decode(&p)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.LoginUser(r.Context(), database.LoginUserParams{
		Email:    p.Email,
		Password: p.Password,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error logging in user: %v", err))
		return
	}

	respondWithJSON(w, 200, databaseUserToUser(user))
}
func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, databaseUserToUser(user))
}
