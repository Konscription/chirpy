package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Konscription/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	type refreshParams struct {
		Token string `json:"token"`
	}

	// check request header for Authorization: Bearer <refreshToken>
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "token not found", err)
		return
	}

	// look up token in the database to get the user id
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil && err == sql.ErrNoRows {
		respondWithError(w, http.StatusUnauthorized, "Failed to get user from refresh token", err)
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database Issue", err)
		return
	}

	// create a new access token
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create access token", err)
		return
	}

	// return the token
	accToken := refreshParams{
		Token: accessToken,
	}
	respondWithJSON(w, http.StatusOK, accToken)
}
