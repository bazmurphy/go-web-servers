package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const (
		port         = "8080"
		filepathRoot = "."
	)

	cfg := &apiConfig{}

	mux := http.NewServeMux()

	handleFiles := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/*", cfg.middlewareMetricsInc(handleFiles))
	mux.HandleFunc("GET /healthz", handleReadiness)
	mux.HandleFunc("GET /metrics", cfg.handleMetrics)
	mux.HandleFunc("GET /reset", cfg.handleReset)

	server := &http.Server{
		Addr: ":" + port,
		// to apply middleware to all routes we wrap the mux in it
		Handler: middlewareLog(mux),
	}

	log.Printf("serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits reset to %d", cfg.fileserverHits)))
}
