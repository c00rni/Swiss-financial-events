package main

import (
	"database/sql"
	"github.com/c00rni/Swiss-financial-events/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"time"
)

type apiConfig struct {
	port string
	DB   *database.Queries
}

func handleReadyness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	return
}

func (cfg *apiConfig) handleGetEvents(w http.ResponseWriter, r *http.Request) {

	location := r.URL.Query().Get("location")
	if location == "" {
		location = "%"
	}
	category := r.URL.Query().Get("category")
	if category == "" {
		category = "%"
	}
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		topic = "%"
	}
	params := database.GetFilteredEventsParams{
		Location: location,
		Name:     category,
		Name_2:   topic,
	}

	events, err := cfg.DB.GetFilteredEvents(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	respondWithJSON(w, http.StatusOK, events)
	return
}

func (cfg *apiConfig) handleGetCategories(w http.ResponseWriter, r *http.Request) {
	names, err := cfg.DB.GetCategories(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	respondWithJSON(w, http.StatusOK, names)
	return
}

func (cfg *apiConfig) handleGetTopics(w http.ResponseWriter, r *http.Request) {
	names, err := cfg.DB.GetTopics(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	respondWithJSON(w, http.StatusOK, names)
	return
}

func main() {
	mux := http.NewServeMux()
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("'PORT' environment variable is not set.")
		return
	}

	dbPath := os.Getenv("SQLI_PATH")
	if dbPath == "" {
		log.Fatalln("'SQLI_PATH' environment variable is not set.")
		return
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

	defer db.Close()
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		port: port,
		DB:   dbQueries,
	}

	mux.HandleFunc("GET /api/healthz", handleReadyness)
	mux.HandleFunc("GET /api/events", apiCfg.handleGetEvents)
	mux.HandleFunc("GET /api/categories", apiCfg.handleGetCategories)
	mux.HandleFunc("GET /api/topics", apiCfg.handleGetTopics)

	go apiCfg.scrapCfasociety(time.Hour * 6)
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + apiCfg.port,
	}

	log.Printf("Serving on port: %s\n", apiCfg.port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Error: %v", err)
	}
}
