package server

import (
	"fmt"
	"net/http"

	organizationHandlers "github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/organization"
	raportHandlers "github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/raport"
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
	router.Post("/auth/signup", userHandlers.RegisterPrivateUser)
	router.Post("/auth/login", userHandlers.LoginUser)
	router.Patch("/me/password", userHandlers.UpdateLoginCredentials)

	// Endpoints that require the JWT
	router.Group(func(protectedRouter chi.Router) {
		protectedRouter.Use(middleware.JWTClaimsMiddleware)

		protectedRouter.Get("/vehicles", vehicleHandlers.GetVehicles)
		protectedRouter.Post("/vehicles", vehicleHandlers.AddVehicle)
		protectedRouter.Get("/vehicles/{id}", vehicleHandlers.GetVehicleById)
		protectedRouter.Patch("/vehicles/{id}", vehicleHandlers.UpdateVehicle)
		protectedRouter.Delete("/vehicles/{id}", vehicleHandlers.DeleteVehicle)

		protectedRouter.Get("/me", userHandlers.GetMyUserInfo)
		protectedRouter.Get("/me/raports", raportHandlers.GetRaports)
		protectedRouter.Get("/me/organization", organizationHandlers.GetMyOrganizationInfo)

		protectedRouter.Patch("/admin/organization", organizationHandlers.PatchMyOrganization)
		protectedRouter.Post("/admin/users", userHandlers.RegisterCorporateUser)

		protectedRouter.Delete("/users/{id}", userHandlers.DeleteUserById)

		protectedRouter.Post("/root/admins", userHandlers.RegisterUserRoot)
		protectedRouter.Post("/root/organizations", organizationHandlers.CreateOrganization)
		protectedRouter.Get("/root/organizations", organizationHandlers.GetAllOrganizations)
		protectedRouter.Delete("/root/organizations/{id}", organizationHandlers.DeleteOrganization)
	})
}
