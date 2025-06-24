package restapi

import (
	"fmt"
	"net/http"

	vehicleHandlers "veh-sense-backend/internal/handlers/rest-api/vehicle"

	"github.com/go-chi/chi"
)

func InitializeAndStart() error {
	router := chi.NewRouter()
	initializeHandlers(router)
	fmt.Println("Starting the HTTP server on port 8080...")
	return http.ListenAndServe(":8080", router)
}

func initializeHandlers(router *chi.Mux) {
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	router.Post("/vehicle", vehicleHandlers.AddVehicle)
}
