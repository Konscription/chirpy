package main

import (
	"encoding/json"
	"net/http"

	"github.com/Konscription/chirpy/internal/auth"
	"github.com/Konscription/chirpy/internal/database"
)

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	type updateParams struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	// check header for access token
	userID, shouldReturn := hasValidAuthToken(r, w, cfg)
	if shouldReturn {
		return
	}

	// Decode the request body
	var params updateParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	hashPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	// Update the user in the database
	err = cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		HashedPassword: hashPass,
		Email:          params.Email})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	user, err := cfg.db.LookupUserById(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get user", err)
		return
	}
	returnedUser := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, returnedUser)
}
