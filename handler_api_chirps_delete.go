package main

import (
	"net/http"

	"github.com/Konscription/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// check header for access token
	accesstoken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	// Validate the token
	userID, err := auth.ValidateJWT(accesstoken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	// Extract the chirp ID from the URL
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	// Get the chirp from the database
	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to find chirp", err)
		return
	}

	// Check if the user is authorized to delete the chirp
	if dbChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You are not authorized to delete this chirp", nil)
		return
	}

	// Call the database function to delete the chirp
	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete chirp", err)
		return
	}

	// Respond with a success message
	respondWithJSON(w, http.StatusNoContent, nil)
}
