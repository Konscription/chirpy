package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type email struct {
		Email string `json:"email"`
	}

	if r.Method != http.MethodPost {
		errorM := ErrorResponse{
			Error: "Method not allowed",
		}
		writeJSONResponse(w, 405, errorM)
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := email{}
	err := decoder.Decode(&params)
	if err != nil {
		errorM := ErrorResponse{
			Error: "Something went wrong",
		}
		writeJSONResponse(w, http.StatusInternalServerError, errorM)
		return
	}
	//Create the user in your database
	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		errorM := ErrorResponse{
			Error: "Something went wrong",
		}
		writeJSONResponse(w, http.StatusInternalServerError, errorM)
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
	writeJSONResponse(w, http.StatusCreated, userResponse)
}
