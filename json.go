package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, statusCode int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if statusCode > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, statusCode, errorResponse{
		Error: msg,
	})
}
