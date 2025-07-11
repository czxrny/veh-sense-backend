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
	// Public endpoints
	router.Post("/auth/register", userHandlers.RegisterPrivateUser)
	router.Post("/user/login", userHandlers.LoginUser)
	router.Post("/user/login/update", userHandlers.UpdateLoginCredentials)

	// Endpoints that require the JWT
	router.Group(func(protectedRouter chi.Router) {
		protectedRouter.Use(middleware.JWTClaimsMiddleware)

		protectedRouter.Get("/vehicle", vehicleHandlers.GetVehicles)
		protectedRouter.Post("/vehicle", vehicleHandlers.AddVehicle)
		protectedRouter.Get("/vehicle/{id}", vehicleHandlers.GetVehicleById)
		protectedRouter.Patch("/vehicle/{id}", vehicleHandlers.UpdateVehicle)
		protectedRouter.Delete("/vehicle/{id}", vehicleHandlers.DeleteVehicle)

		protectedRouter.Post("/admin/user", userHandlers.RegisterCorporateUser)
		protectedRouter.Delete("/user/{id}", userHandlers.DeleteUserById)

		// For organizations
		protectedRouter.Get("/me/organization", organizationHandlers.GetMyOrganizationInfo)
		protectedRouter.Patch("/me/organization", organizationHandlers.PatchMyOrganization)

		// Root only
		protectedRouter.Post("/root/user", userHandlers.RegisterUserRoot)

		// Organization related
		protectedRouter.Post("/root/organization", organizationHandlers.CreateOrganization)
		protectedRouter.Get("/root/organization", organizationHandlers.GetAllOrganizations)
		protectedRouter.Delete("/root/organization/{id}", organizationHandlers.DeleteOrganization)
	})
}
