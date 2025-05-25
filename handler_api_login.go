package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Konscription/chirpy/internal/auth"
	"github.com/Konscription/chirpy/internal/database"
)

const refreshTokenExpirationTime = time.Hour * 24 * 60 // 60 days

// handler function to login a user
func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	// Decode the request body
	type loginParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	var params loginParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Validate the email and password exist
	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}

	// Authenticate the user
	user, err := cfg.db.LookupUserbyEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	// Check if the password is correct
	if err := auth.CheckPasswordHash(user.HashedPassword, params.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// Generate a JWT token
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return
	}
	// Generate a refresh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Failed to generate refresh token: %v", err)
		return
	}

	refreshTokenParams := database.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(refreshTokenExpirationTime),
	}

	// Store the refresh token in the database
	_, err = cfg.db.CreateRefreshToken(r.Context(), refreshTokenParams)
	if err != nil {
		log.Printf("Failed to store refresh token: %v", err)
		return
	}
	// Respond with the user information (excluding password)
	resp := response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refreshToken,
	}
	respondWithJSON(w, http.StatusOK, resp)
}
