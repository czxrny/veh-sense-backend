package server

import (
	"fmt"
	"net/http"

	database "github.com/czxrny/veh-sense-backend/rest-api/internal/app"
	organizationHandlers "github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/organization"
	raportHandlers "github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/raport"
	userHandlers "github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/user"
	vehicleHandlers "github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/vehicle"
	"github.com/czxrny/veh-sense-backend/rest-api/internal/middleware"
	"github.com/go-chi/chi"
)

func InitializeAndStart(app *database.App) error {
	router := initializeHandlers(app)
	fmt.Println("Starting the HTTP server on port 8080...")
	return http.ListenAndServe(":8080", router)
}

func initializeHandlers(app *database.App) *chi.Mux {
	vehHandler := vehicleHandlers.NewVehicleHandler(app.VehicleService)

	router := chi.NewRouter()
	// Public endpoints
	router.Post("/auth/signup", userHandlers.RegisterPrivateUser)
	router.Post("/auth/login", userHandlers.LoginUser)
	router.Patch("/me/password", userHandlers.UpdateLoginCredentials)

	// Endpoints that require the JWT
	router.Group(func(protectedRouter chi.Router) {
		protectedRouter.Use(middleware.JWTClaimsMiddleware)

		protectedRouter.Get("/vehicles", vehHandler.GetVehicles)
		protectedRouter.Post("/vehicles", vehHandler.AddVehicle)
		protectedRouter.Get("/vehicles/{id}", vehHandler.GetVehicleById)
		protectedRouter.Patch("/vehicles/{id}", vehHandler.UpdateVehicle)
		protectedRouter.Delete("/vehicles/{id}", vehHandler.DeleteVehicle)

		protectedRouter.Get("/raports", raportHandlers.GetRaports)
		protectedRouter.Delete("/raports/{id}", raportHandlers.DeleteRaport)

		protectedRouter.Get("/me", userHandlers.GetMyUserInfo)
		protectedRouter.Get("/me/organization", organizationHandlers.GetMyOrganizationInfo)

		protectedRouter.Patch("/admin/organization", organizationHandlers.PatchMyOrganization)
		protectedRouter.Post("/admin/users", userHandlers.RegisterCorporateUser)
		protectedRouter.Get("/admin/users", userHandlers.GetAllUsersInfo)

		protectedRouter.Delete("/users/{id}", userHandlers.DeleteUserById)

		protectedRouter.Post("/root/admins", userHandlers.RegisterUserRoot)
		protectedRouter.Post("/root/organizations", organizationHandlers.CreateOrganization)
		protectedRouter.Get("/root/organizations", organizationHandlers.GetAllOrganizations)
		protectedRouter.Delete("/root/organizations/{id}", organizationHandlers.DeleteOrganization)
	})

	return router
}
