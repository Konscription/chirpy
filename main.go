package main

import (
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
	mux.HandleFunc("GET /api/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /api/reset", apiCfg.resetHandler)
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
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// get the count and format it correctcly
	hitCount := cfg.fileserverHits.Load()
	w.Write([]byte(fmt.Sprintf("Hits: %d", hitCount)))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Reset the hit counter
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Counter reset"))
}
