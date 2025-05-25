package main

import (
	"net/http"

	"github.com/Konscription/chirpy/internal/auth"
	"github.com/google/uuid"
)

func hasValidAuthToken(r *http.Request, w http.ResponseWriter, cfg *apiConfig) (userID uuid.UUID, shouldReturn bool) {
	accesstoken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return uuid.UUID{}, true
	}
	// Validate the token
	userID, err = auth.ValidateJWT(accesstoken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return uuid.UUID{}, true
	}
	return userID, false
}
