package server

import (
	"fmt"
	"net/http"

	organizationHandlers "github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/organization"
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

	router.Post("/auth/register", userHandlers.RegisterPrivateUser)
	router.Post("/admin/user", userHandlers.RegisterCorporateUser)
	router.Post("/root/user", userHandlers.RegisterUserRoot)
	router.Post("/user/login", userHandlers.LoginUser)
	router.Delete("/user/{id}", userHandlers.DeleteUserById)

	// Organization related - root only
	router.Post("/root/organization", organizationHandlers.CreateOrganization)
	router.Get("/root/organization", organizationHandlers.GetAllOrganizations)
	router.Delete("/root/organization", organizationHandlers.DeleteOrganization)

	// For organizations
	router.Get("/me/organization", organizationHandlers.GetMyOrganizationInfo)
	router.Get("/root/organization", organizationHandlers.PatchMyOrganization)
}
