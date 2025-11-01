package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/Tom-Webbo/Go-HTTP-Server/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

type User struct {
	ID         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Email      string    `json:"email"`
}

type chirpRecord struct {
	ID         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Body       string    `json:"body"`
	UserID     uuid.UUID `json:"user_id"`
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	envType := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to SQL - %s", err)
	}
	dbQueries := database.New(db)

	acfg := apiConfig{db: dbQueries, platform: envType}
	router := http.NewServeMux()

	router.Handle("/app/", acfg.middlewareHandleIncrementMetric(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	router.HandleFunc("GET /api/healthz", handleHealthz)
	router.HandleFunc("GET /admin/metrics", acfg.handleMetricRetrieval)
	router.HandleFunc("POST /admin/reset", acfg.handleReset)
	router.HandleFunc("POST /api/chirps", acfg.handleCreateChirp)
	router.HandleFunc("POST /api/users", acfg.handleUserCreation)
	router.HandleFunc("GET /api/chirps", acfg.handleGetChirps)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("server failure - %v", err)
	}

}
