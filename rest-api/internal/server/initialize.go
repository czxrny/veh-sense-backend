package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	internal "github.com/czxrny/veh-sense-backend/rest-api/internal/app"
	o "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/organization/handler"
	r "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/report/handler"
	u "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/user/handler"
	v "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/vehicle/handler"
	"github.com/czxrny/veh-sense-backend/shared/middleware"
	"github.com/go-chi/chi"
)

var AppStart = time.Now()

func InitializeAndStart(app *internal.App) error {
	router := initializeHandlers(app)
	port := os.Getenv("REST_API_PORT")
	fmt.Printf("Starting the HTTP server on port %s...\n", port)
	return http.ListenAndServe(":"+port, router)
}

func initializeHandlers(app *internal.App) *chi.Mux {
	vehHandler := v.NewVehicleHandler(app.VehicleService)
	orgHandler := o.NewOrganizationHandler(app.OrganizationService)
	rapHandler := r.NewReportHandler(app.ReportService)
	userAuthHandler := u.NewUserAuthHandler(app.UserService)
	userInfoHandler := u.NewUserInfoHandler(app.UserService)

	router := chi.NewRouter()
	// Public endpoints
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong"))
	})
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"Function":   "Rest Api",
			"Started at": AppStart.Format("02-01-2006 15:04:05 MST"),
			"Uptime":     time.Since(AppStart).String(),
		})
	})

	router.Group(func(authRouter chi.Router) {
		authRouter.Use(middleware.CORSMiddleware)
		authRouter.Use(middleware.RequireAPIKeyMiddleware)

		authRouter.Post("/auth/signup", userAuthHandler.RegisterPrivateUser)
		authRouter.Post("/auth/login", userAuthHandler.LoginUser)
		authRouter.Patch("/me/credentials", userAuthHandler.UpdateLoginCredentials)
		authRouter.Post("/auth/refresh", userAuthHandler.RefreshByKey)
	})

	// Endpoints that require the JWT
	router.Group(func(protectedRouter chi.Router) {
		protectedRouter.Use(middleware.CORSMiddleware)
		protectedRouter.Use(middleware.JWTClaimsMiddleware)

		protectedRouter.Get("/vehicles", vehHandler.GetVehicles)
		protectedRouter.Post("/vehicles", vehHandler.AddVehicle)
		protectedRouter.Get("/vehicles/{id}", vehHandler.GetVehicleById)
		protectedRouter.Patch("/vehicles/{id}", vehHandler.UpdateVehicle)
		protectedRouter.Delete("/vehicles/{id}", vehHandler.DeleteVehicle)

		protectedRouter.Get("/reports", rapHandler.GetReports)
		protectedRouter.Get("/reports/{id}/data", rapHandler.GetReportDataById)
		protectedRouter.Get("/reports/{id}", rapHandler.GetReportById)
		protectedRouter.Delete("/reports/{id}", rapHandler.DeleteReport)

		protectedRouter.Get("/me", userInfoHandler.GetMyUserInfo)
		protectedRouter.Get("/me/organization", orgHandler.GetMyOrganizationInfo)

		protectedRouter.Patch("/admin/organization", orgHandler.UpdateMyOrganization)
		protectedRouter.Post("/admin/users", userAuthHandler.RegisterCorporateUser)
		protectedRouter.Get("/admin/users", userInfoHandler.GetAllUsersInfo)
		protectedRouter.Get("/admin/reports", rapHandler.GetReportsAdmin)

		protectedRouter.Delete("/users/{id}", userInfoHandler.DeleteUserById)
		protectedRouter.Get("/users/{id}", userInfoHandler.GetUserInfoById)

		protectedRouter.Post("/root/admins", userAuthHandler.RegisterUserRoot)
		protectedRouter.Post("/root/organizations", orgHandler.CreateOrganization)
		protectedRouter.Get("/root/organizations", orgHandler.GetAllOrganizations)
		protectedRouter.Delete("/root/organizations/{id}", orgHandler.DeleteOrganization)
	})

	return router
}
