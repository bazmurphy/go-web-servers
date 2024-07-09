package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const (
		port         = "8080"
		filepathRoot = "."
	)

	cfg := &apiConfig{
		fileserverHits: 0,
	}

	mux := http.NewServeMux()

	handlerFileserver := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/*", cfg.middlewareMetricsInc(handlerFileserver))

	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("GET /api/admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("GET /api/admin/reset", cfg.handlerReset)

	server := &http.Server{
		Addr: ":" + port,
		// to apply middleware to all routes we wrap the mux in it
		Handler: middlewareLog(mux),
	}

	log.Printf("serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

type Chirp struct {
	Body string `json:"body"`
}

type Response struct {
	Error string `json:"error,omitempty"` // (!) the omitempty tag tells the JSON encoder to omit the field if it's empty
	Valid bool   `json:"valid,omitempty"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		response := Response{
			Error: "Something went wrong",
		}
		jsonResponseBody, err := json.Marshal(response)
		if err != nil {
			log.Printf("error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		log.Println(string(jsonResponseBody))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(jsonResponseBody)
		return
	}

	if len(chirp.Body) > 140 {
		response := Response{
			Error: "Chirp is too long",
		}
		jsonResponseBody, err := json.Marshal(response)
		if err != nil {
			log.Printf("error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		log.Println(string(jsonResponseBody))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(jsonResponseBody)
		return
	}

	response := Response{
		Valid: true,
	}
	jsonResponseBody, err := json.Marshal(response)
	if err != nil {
		log.Printf("error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	log.Println(string(jsonResponseBody))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonResponseBody)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html>

<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>

</html>
`, cfg.fileserverHits)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits reset to %d", cfg.fileserverHits)))
}
