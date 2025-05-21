package main

import (
	"net/http"
)

// Handler function to get chirps
func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	// Get the chirps from the database
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get chirps", err)
		return
	}
	chirpsResponse := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		chirpsResponse[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	// Respond with the chirps
	respondWithJSON(w, http.StatusOK, chirpsResponse)
}
