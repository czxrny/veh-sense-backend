package server

import (
	"fmt"
	"net/http"

	vehicleHandlers "github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/vehicle"

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
	router.Get("/vehicle", vehicleHandlers.GetVehicles)
	router.Post("/vehicle", vehicleHandlers.AddVehicle)
	router.Get("/vehicle/{id}", vehicleHandlers.GetVehicleById)
	router.Patch("/vehicle/{id}", vehicleHandlers.UpdateVehicle)
	router.Delete("/vehicle/{id}", vehicleHandlers.DeleteVehicle)
}
