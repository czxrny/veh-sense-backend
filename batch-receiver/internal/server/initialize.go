package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

var appStart = time.Now()

func InitializeAndStart() error {
	port := "8090"
	router := initializeHandlers()
	fmt.Printf("Starting the HTTP BATCH Receiver server on port %s...\n", port)
	return http.ListenAndServe(":"+port, router)
}

func initializeHandlers() *chi.Mux {
	router := chi.NewRouter()
	// Public endpoints
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong"))
	})
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"Function":   "Batch Receiver",
			"Started at": appStart.Format("02-01-2006 15:04:05 MST"),
			"Uptime":     time.Since(appStart).String(),
		})
	})

	return router
}
