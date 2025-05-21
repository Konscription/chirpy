package main

import (
	"net/http"

	"github.com/google/uuid"
)

// handler function to get a chirp
func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	// Get the chirp ID from the URL
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		respondWithError(w, http.StatusNotFound, "Chirp ID is required", nil)
		return
	}
	// convert chirpID to UUID
	validChirpID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid chirp ID", err)
		return
	}

	// Get the chirp from the database
	dbChirp, err := cfg.db.GetChirp(r.Context(), validChirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to find chirp", err)
		return
	}

	// Convert the database chirp to the API chirp
	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	// Respond with the chirp
	respondWithJSON(w, http.StatusOK, chirp)
}
