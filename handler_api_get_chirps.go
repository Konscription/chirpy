package main

import (
	"net/http"

	"github.com/Konscription/chirpy/internal/database"
	"github.com/google/uuid"
)

// Handler function to get chirps
func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	authorID := r.URL.Query().Get("author_id")
	var (
		dbChirps []database.Chirp
		err      error
	)

	if authorID != "" {
		authorUUID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id format", err)
			return
		}
		dbChirps, err = cfg.db.GetChirpsByAuthor(r.Context(), authorUUID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to get chirps for author", err)
			return
		}
	} else {
		dbChirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to get chirps", err)
			return
		}
	}

	chirpsResponse := make([]Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		chirpsResponse[i] = Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, chirpsResponse)
}
