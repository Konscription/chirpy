package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	const filepathroot = "."
	const port = "8080"

	// First, create an instance of apiConfig
	apiCfg := &apiConfig{}
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(filepathroot))

	// Then use the instance of call the middleware method
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /api/healthz", handlerRediness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.validateChirpHandler)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("serving files from %s on port: %s\n", filepathroot, port)
	log.Fatal(server.ListenAndServe())
}

func handlerRediness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// get the count and format it correctcly
	hitCount := cfg.fileserverHits.Load()
	w.Write([]byte(fmt.Sprintf(`<html>
  <body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
  </body>
</html>`, hitCount)))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Reset the hit counter
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Counter reset"))
}

func (cfg *apiConfig) validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type ErrorResponse struct {
		Error string `json:"error"`
	}

	type ValidResponse struct {
		Valid bool `json:"valid"`
	}

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
	validM := ValidResponse{
		Valid: true,
	}
	log.Printf("Chirp is valid")
	writeJSONResponse(w, 200, validM)
}

func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(response)
}
