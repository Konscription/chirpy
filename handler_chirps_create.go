package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Konscription/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	// Check if the request method is POST
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	// Validate the chirp
	if len(params.Body) > 140 {
		errorM := ErrorResponse{
			Error: "Chirp is too long",
		}
		log.Printf("Chirp is too long")
		respondWithJSON(w, http.StatusBadRequest, errorM)
		return
	}

	// clean up chirp if it contains any profane words
	CleanedBody, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	// create a database chirp object
	validChirp := database.CreateChirpParams{
		Body:   CleanedBody,
		UserID: params.UserID,
	}
	// insert the chirp into the database
	createdChirp, err := cfg.db.CreateChirp(r.Context(), validChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body:      createdChirp.Body,
		UserID:    createdChirp.UserID,
	})
}

func chirpCleaner(chirp string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	// check if the chirp contains any profane words
	splitChirp := strings.Split(chirp, " ")
	cleanChirp := []string{}
	// replace the profane word with asterisks
	for _, word := range splitChirp {
		if slices.Contains(profaneWords, strings.ToLower(word)) {
			cleanChirp = append(cleanChirp, "****")
			continue
		}
		cleanChirp = append(cleanChirp, word)
	}
	return strings.Join(cleanChirp, " ")
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	cleanedBody := chirpCleaner(body)
	return cleanedBody, nil
}
