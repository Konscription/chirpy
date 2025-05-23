package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Konscription/chirpy/internal/auth"
	"github.com/Konscription/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if r.Method != http.MethodPost {
		errorM := ErrorResponse{
			Error: "Method not allowed",
		}
		respondWithJSON(w, 405, errorM)
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorM := ErrorResponse{
			Error: "Something went wrong",
		}
		respondWithJSON(w, http.StatusInternalServerError, errorM)
		return
	}
	//hash the password
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		errorM := ErrorResponse{
			Error: "issue with password hashing",
		}
		respondWithJSON(w, http.StatusInternalServerError, errorM)
		return
	}

	//Create the user in your database
	dbParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}
	user, err := cfg.db.CreateUser(r.Context(), dbParams)
	if err != nil {
		errorM := ErrorResponse{
			Error: "Something went wrong",
		}
		respondWithJSON(w, http.StatusInternalServerError, errorM)
		return
	}
	//Map the returned database user to your own User struct
	userResponse := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	//Write the response back as JSON
	respondWithJSON(w, http.StatusCreated, userResponse)
}
