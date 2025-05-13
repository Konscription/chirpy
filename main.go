package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathroot = "."
	const port = "8080"

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(filepathroot))
	mux.Handle("/app/", http.StripPrefix("/app", fileServer))
	mux.HandleFunc("/healthz", handlerRediness)

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
