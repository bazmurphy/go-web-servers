package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
)

type apiConfig struct {
	fileserverHits int
	db             *DB
}

func main() {
	const (
		port         = "8080"
		filepathRoot = "."
	)

	db, err := NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	cfg := &apiConfig{
		fileserverHits: 0,
		db:             db,
	}

	mux := http.NewServeMux()

	handlerFileserver := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/*", cfg.middlewareMetricsInc(handlerFileserver))

	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	mux.HandleFunc("POST /api/chirps", cfg.handlerChirpsPost)
	mux.HandleFunc("GET /api/chirps", cfg.handlerChirpsGet)

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
	ID   int    `json:"id,omitempty"`
	Body string `json:"body"`
}

type Response struct {
	// (!) the omitempty tag tells the JSON encoder to omit the field if it's empty
	Error       string `json:"error,omitempty"`
	Valid       bool   `json:"valid,omitempty"`
	CleanedBody string `json:"cleaned_body,omitempty"`
	ID          int    `json:"id,omitempty"`
	Body        string `json:"body,omitempty"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		response := Response{
			Error: "Something went wrong",
		}
		byteData, err := json.Marshal(response)
		if err != nil {
			log.Printf("error marshalling JSON: %s", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(byteData)
		return
	}

	if len(chirp.Body) > 140 {
		response := Response{
			Error: "Chirp is too long",
		}
		byteData, err := json.Marshal(response)
		if err != nil {
			log.Printf("error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(byteData)
		return
	}

	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	wordsSplit := strings.Split(chirp.Body, " ")
	for index, word := range wordsSplit {
		for _, profaneWord := range profaneWords {
			if strings.ToLower(word) == profaneWord {
				wordsSplit[index] = "****"
			}
		}
	}
	wordsRejoined := strings.Join(wordsSplit, " ")

	if chirp.Body != wordsRejoined {
		response := Response{
			CleanedBody: wordsRejoined,
		}
		byteData, err := json.Marshal(response)
		if err != nil {
			log.Printf("error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(byteData)
		return
	}

	response := Response{
		Valid: true,
	}
	byteData, err := json.Marshal(response)
	if err != nil {
		log.Printf("error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(byteData)
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

func (cfg *apiConfig) handlerChirpsPost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		response := Response{
			Error: "Something went wrong",
		}
		byteData, err := json.Marshal(response)
		if err != nil {
			log.Printf("error marshalling JSON: %s", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(byteData)
		return
	}

	if len(chirp.Body) > 140 {
		response := Response{
			Error: "Chirp is too long",
		}
		byteData, err := json.Marshal(response)
		if err != nil {
			log.Printf("error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(byteData)
		return
	}

	chirp, err = cfg.db.CreateChirp(chirp.Body)
	if err != nil {
		// TODO handle this error
		w.WriteHeader(500)
		return
	}

	byteData, err := json.Marshal(chirp)
	if err != nil {
		log.Printf("error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(byteData)
}

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		// TODO handle error
		w.WriteHeader(500)
		return
	}

	byteData, err := json.Marshal(chirps)
	if err != nil {
		log.Printf("error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(byteData)
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	dbOnDisk := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	err := dbOnDisk.ensureDB()
	if err != nil {
		return nil, err
	}
	return dbOnDisk, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbMemory, err := db.loadDB()
	if err != nil {
		return Chirp{}, nil
	}
	chirp := Chirp{
		ID:   len(dbMemory.Chirps) + 1,
		Body: body,
	}
	dbMemory.Chirps[chirp.ID] = chirp
	err = db.writeDB(dbMemory)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbMemory, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := make([]Chirp, 0, len(dbMemory.Chirps))
	for _, chirp := range dbMemory.Chirps {
		chirps = append(chirps, chirp)
	}
	sort.Slice(chirps, func(a, b int) bool {
		return chirps[a].ID < chirps[b].ID
	})
	return chirps, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbExistingFile, err := os.Open(db.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			dbNewFile, err := os.Create(db.path)
			if err != nil {
				return err
			}
			defer dbNewFile.Close()
			_, err = dbNewFile.Write([]byte("{}")) // (!) an empty json object
			if err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}
	defer dbExistingFile.Close()
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	byteData, err := os.ReadFile(db.path)
	if err != nil {
		log.Println(err)
		return DBStructure{}, nil
	}

	dbInMemory := &DBStructure{
		Chirps: make(map[int]Chirp),
	}

	err = json.Unmarshal(byteData, dbInMemory)
	if err != nil {
		return DBStructure{}, err
	}

	return *dbInMemory, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	json, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, json, 0644)
	if err != nil {
		return err
	}

	return nil
}
