package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Konscription/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func main() {
	const filepathroot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}
	// Open a connection to the database
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer dbConn.Close()
	dbQueries := database.New(dbConn)

	// First, create an instance of apiConfig
	apiCfg := &apiConfig{}
	apiCfg.db = dbQueries

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(filepathroot))

	// Then use the instance of call the middleware method
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))

	mux.HandleFunc("GET /api/healthz", handlerRediness)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.validateChirpHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
