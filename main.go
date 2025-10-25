package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {

	acfg := apiConfig{}
	router := http.NewServeMux()

	router.Handle("/app/", acfg.middlewareHandleIncrementMetric(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	router.HandleFunc("GET /api/healthz", handleHealthz)
	router.HandleFunc("GET /admin/metrics", acfg.handleMetricRetrieval)
	router.HandleFunc("POST /admin/reset", acfg.handleResetMetrics)
	router.HandleFunc("POST /api/validate_chirp", acfg.handleValidateChirps)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("server failure - %v", err)
	}

}
