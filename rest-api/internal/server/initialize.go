package server

import (
	"fmt"
	"net/http"

	userHandlers "github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/user"
	vehicleHandlers "github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/vehicle"
	"github.com/czxrny/veh-sense-backend/rest-api/internal/middleware"
	"github.com/go-chi/chi"
)

func InitializeAndStart() error {
	router := chi.NewRouter()
	initializeHandlers(router)
	fmt.Println("Starting the HTTP server on port 8080...")
	return http.ListenAndServe(":8080", router)
}

func initializeHandlers(router *chi.Mux) {
	router.Use(middleware.JWTClaimsMiddleware)

	router.Get("/vehicle", vehicleHandlers.GetVehicles)
	router.Post("/vehicle", vehicleHandlers.AddVehicle)
	router.Get("/vehicle/{id}", vehicleHandlers.GetVehicleById)
	router.Patch("/vehicle/{id}", vehicleHandlers.UpdateVehicle)
	router.Delete("/vehicle/{id}", vehicleHandlers.DeleteVehicle)

	router.Post("/user/register", userHandlers.RegisterUser)
	router.Post("/user/login", userHandlers.LoginUser)
	router.Delete("/user/{id}", userHandlers.DeleteUserById)

}
