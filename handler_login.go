package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Konscription/chirpy/internal/auth"
)

const defaultExpirationTime = 3600 // 1 hour

// handler function to login a user
func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	// Decode the request body
	type loginParams struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds,omitempty"`
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

	// Create exparation time
	// Check if the expires_in_seconds parameter is provided
	// If not, set a default expiration time
	var expiration time.Duration
	if params.ExpiresInSeconds != nil && *params.ExpiresInSeconds <= defaultExpirationTime {
		expiration = time.Duration(*params.ExpiresInSeconds) * time.Second
	} else {
		// Default expiration time
		expiration = defaultExpirationTime * time.Second // 1 hour
	}
	// Generate a JWT token
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiration)
	if err != nil {
		// Log the error but do not expose internal details to the client
		log.Printf("Failed to generate token: %v", err)
		//respondWithError(w, http.StatusInternalServerError, "Failed to generate token", nil)
		return
	}
	// Respond with the user information (excluding password)
	response := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	}
	respondWithJSON(w, http.StatusOK, response)
}
