package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type apiConfig struct {
	port string
}

func handleReadyness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
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

	apiCfg := apiConfig{
		port: port,
	}

	mux.HandleFunc("GET /api/healthz", handleReadyness)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + apiCfg.port,
	}

	log.Printf("Serving on port: %s\n", apiCfg.port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("Error: %v", err)
	}
}
