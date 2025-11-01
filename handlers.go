package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Tom-Webbo/Go-HTTP-Server/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func (cfg *apiConfig) middlewareHandleIncrementMetric(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetricRetrieval(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Only available in dev enviroments", nil)
		return
	}
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to delete all users", err)
		return
	}
	respondWithJson(w, http.StatusOK, nil)
}

// func (cfg *apiConfig) handleValidateChirps(w http.ResponseWriter, r *http.Request) {

// 	type parameters struct {
// 		Body string `json:"body"`
// 	}

// 	type returnVals struct {
// 		CleanedBody string `json:"cleaned_body"`
// 	}

// 	params := parameters{}
// 	decoder := json.NewDecoder(r.Body)

// 	if err := decoder.Decode(&params); err != nil {
// 		respondWithError(w, http.StatusBadRequest, "Something went wrong", err)
// 		return
// 	}
// 	if len(params.Body) > 140 {
// 		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
// 		return
// 	}

// 	msg := cleaner(params.Body)

// 	respondWithJson(w, http.StatusOK, returnVals{CleanedBody: msg})

// }

func (cfg *apiConfig) handleUserCreation(w http.ResponseWriter, r *http.Request) {

	type body struct {
		Email string `json:"email"`
	}

	params := body{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to decode body", err)
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to create new user", err)
	}

	respondWithJson(w, http.StatusCreated, User{
		ID:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	})

}

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "fails to decode json", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	msg := cleaner(params.Body)
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: msg, UserID: params.UserID})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to create user", err)
		return
	}

	respondWithJson(w, http.StatusCreated, chirpRecord{
		ID:         chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body:       chirp.Body,
		UserID:     chirp.UserID,
	})

}

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	chirpSlice := []chirpRecord{}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to retrieve chirps", err)
	}
	for _, chirp := range chirps {
		chirpSlice = append(chirpSlice, chirpRecord{ID: chirp.ID, Created_at: chirp.CreatedAt, Updated_at: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID})
	}
	type returnVals struct {
		Body []chirpRecord
	}
	respondWithJson(w, http.StatusOK, chirpSlice)
}
