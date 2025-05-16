package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}
type ValidResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

func (cfg *apiConfig) validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		errorM := ErrorResponse{
			Error: "Method not allowed",
		}
		writeJSONResponse(w, 405, errorM)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorM := ErrorResponse{
			Error: "Something went wrong"}
		log.Printf("Error decoding JSON: %v", err)
		writeJSONResponse(w, 500, errorM)
		return
	}
	// Validate the chirp
	if len(params.Body) > 140 {
		errorM := ErrorResponse{
			Error: "Chirp is too long",
		}
		log.Printf("Chirp is too long")
		writeJSONResponse(w, 400, errorM)
		return
	}

	// check if the chirp contains any profane words
	validM := ValidResponse{
		CleanedBody: chirpProfaneChecker(params.Body),
	}
	log.Printf("Chirp is valid")
	writeJSONResponse(w, 200, validM)
}

func chirpProfaneChecker(chirp string) string {
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
