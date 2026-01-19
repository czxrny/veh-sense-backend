package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	internal "github.com/czxrny/veh-sense-backend/batch-receiver/internal/app"
	r "github.com/czxrny/veh-sense-backend/batch-receiver/internal/domain/raport/handler"
	"github.com/czxrny/veh-sense-backend/shared/middleware"

	"github.com/go-chi/chi/v5"
)

var appStart = time.Now()

func InitializeAndStart(app *internal.App) error {
	router := initializeHandlers(app)
	port := os.Getenv("BATCH_RECEIVER_PORT")
	fmt.Printf("Starting the HTTP BATCH Receiver server on port %s...\n", port)
	return http.ListenAndServe(":"+port, router)
}

func initializeHandlers(app *internal.App) *chi.Mux {
	router := chi.NewRouter()

	uh := r.NewUploadHandler(&app.Service)

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

	router.Group(func(protectedRouter chi.Router) {
		protectedRouter.Use(middleware.JWTClaimsMiddleware)
		protectedRouter.Post("/upload", uh.Upload)
	})

	return router
}
