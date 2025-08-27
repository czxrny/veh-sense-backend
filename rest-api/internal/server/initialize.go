package server

import (
	"fmt"
	"net/http"

	database "github.com/czxrny/veh-sense-backend/rest-api/internal/app"
	o "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/organization/handler"
	r "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/raport/handler"
	u "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/user/handler"
	v "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/vehicle/handler"
	"github.com/czxrny/veh-sense-backend/rest-api/internal/middleware"
	"github.com/go-chi/chi"
)

func InitializeAndStart(app *database.App) error {
	router := initializeHandlers(app)
	fmt.Println("Starting the HTTP server on port 8080...")
	return http.ListenAndServe(":8080", router)
}

func initializeHandlers(app *database.App) *chi.Mux {
	vehHandler := v.NewVehicleHandler(app.VehicleService)
	orgHandler := o.NewOrganizationHandler(app.OrganizationService)
	rapHandler := r.NewRaportHandler(app.RaportService)
	userAuthHandler := u.NewUserAuthHandler(app.UserService)
	userInfoHandler := u.NewUserInfoHandler(app.UserService)

	router := chi.NewRouter()
	// Public endpoints
	router.Post("/auth/signup", userAuthHandler.RegisterPrivateUser)
	router.Post("/auth/login", userAuthHandler.LoginUser)
	router.Patch("/me/credentials", userAuthHandler.UpdateLoginCredentials)

	// Endpoints that require the JWT
	router.Group(func(protectedRouter chi.Router) {
		protectedRouter.Use(middleware.JWTClaimsMiddleware)

		protectedRouter.Get("/vehicles", vehHandler.GetVehicles)
		protectedRouter.Post("/vehicles", vehHandler.AddVehicle)
		protectedRouter.Get("/vehicles/{id}", vehHandler.GetVehicleById)
		protectedRouter.Patch("/vehicles/{id}", vehHandler.UpdateVehicle)
		protectedRouter.Delete("/vehicles/{id}", vehHandler.DeleteVehicle)

		protectedRouter.Get("/raports", rapHandler.GetRaports)
		protectedRouter.Delete("/raports/{id}", rapHandler.DeleteRaport)

		protectedRouter.Get("/me", userInfoHandler.GetMyUserInfo)
		protectedRouter.Get("/me/organization", orgHandler.GetMyOrganizationInfo)

		protectedRouter.Patch("/admin/organization", orgHandler.UpdateMyOrganization)
		protectedRouter.Post("/admin/users", userAuthHandler.RegisterCorporateUser)
		protectedRouter.Get("/admin/users", userInfoHandler.GetAllUsersInfo)

		protectedRouter.Delete("/users/{id}", userInfoHandler.DeleteUserById)

		protectedRouter.Post("/root/admins", userAuthHandler.RegisterUserRoot)
		protectedRouter.Post("/root/organizations", orgHandler.CreateOrganization)
		protectedRouter.Get("/root/organizations", orgHandler.GetAllOrganizations)
		protectedRouter.Delete("/root/organizations/{id}", orgHandler.DeleteOrganization)
	})

	return router
}
