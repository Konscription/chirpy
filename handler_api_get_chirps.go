package main

import (
	"net/http"
	"sort"

	"github.com/Konscription/chirpy/internal/database"
	"github.com/google/uuid"
)

// Handler function to get chirps
func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}
	// get author parameter if it exists
	authorID := r.URL.Query().Get("author_id")
	var (
		dbChirps []database.Chirp
		err      error
	)

	// get sort parameter if it exists
	sortParam := r.URL.Query().Get("sort")
	if sortParam != "" && sortParam != "asc" && sortParam != "desc" {
		respondWithError(w, http.StatusBadRequest, "Invalid sort format. valid values are (asc,desc)", nil)
		return
	} else if sortParam == "" {
		sortParam = "asc"
	}

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
	// sort the slice if desc, as default return from database is asc
	if sortParam == "desc" {
		// Sort chirps in descending order based on the CreatedAt timestamp
		sort.Slice(chirpsResponse, func(i, j int) bool {
			return chirpsResponse[i].CreatedAt.After(chirpsResponse[j].CreatedAt)
		})
	}
	respondWithJSON(w, http.StatusOK, chirpsResponse)
}
