package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	// Check if the request method is POST
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}
	// Decode the JSON request body
	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// is the event a valid one?
	if req.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	// update the user in the database
	err = cfg.db.UpdateUserToChirpyRed(r.Context(), req.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to update user", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
