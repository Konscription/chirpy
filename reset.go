package main

import "net/http"

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Reset the hit counter
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Counter reset"))
}
