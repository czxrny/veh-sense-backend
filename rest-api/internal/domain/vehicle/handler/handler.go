package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	common "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/common/handler"
	s "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/vehicle/service"
	"github.com/czxrny/veh-sense-backend/rest-api/internal/middleware"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

type VehicleHandler struct {
	*s.VehicleService
}

func NewVehicleHandler(vehicleService *s.VehicleService) *VehicleHandler {
	return &VehicleHandler{
		VehicleService: vehicleService,
	}
}

// Returns all the vehicles that are available to user:
// If the user is a private user: returns all of the private cars
// If the user is corporate and is not an admin: returns all of the assigned cars + shared (the ones without the owner_id)
// If the user is an admin of organization: returns all of the organization cars
// Root gets all the info.
func (v *VehicleHandler) GetVehicles(w http.ResponseWriter, r *http.Request) {
	common.GetAllHandler(w, r, func(ctx context.Context, query url.Values) ([]models.Vehicle, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		filter := models.VehicleFilter{
			Brand:          query.Get("brand"),
			MinCapacity:    query.Get("minEngineCapacity"),
			MaxCapacity:    query.Get("maxEngineCapacity"),
			MinEnginePower: query.Get("minEnginePower"),
			MaxEnginePower: query.Get("maxEnginePower"),
			UserID:         authClaims.UserID,
			OrganizationID: authClaims.OrganizationID,
			Role:           authClaims.Role,
		}

		return v.VehicleService.FindVehicles(r.Context(), filter)
	})
}

func (v *VehicleHandler) AddVehicle(w http.ResponseWriter, r *http.Request) {
	common.PostHandler(w, r, func(ctx context.Context, vehicle *models.Vehicle) (*models.Vehicle, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		return v.VehicleService.AddVehicle(ctx, vehicle, authClaims)
	})
}

// Returns the vehicle - only if the user is the owner of the vehicle / the vehicle is shared in the corporation the user is in
func (v *VehicleHandler) GetVehicleById(w http.ResponseWriter, r *http.Request) {
	common.GetByIdHandler(w, r, func(ctx context.Context, id int) (*models.Vehicle, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		return v.VehicleService.GetById(ctx, authClaims, id)
	})
}

func (v *VehicleHandler) UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	common.PatchHandler(w, r, func(ctx context.Context, updatedVehicle *models.VehicleUpdate, id int) (*models.Vehicle, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		return v.VehicleService.UpdateById(ctx, authClaims, updatedVehicle, id)
	})
}

func (v *VehicleHandler) DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	common.DeleteHandler(w, r, func(ctx context.Context, id int) error {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		return v.VehicleService.DeleteById(ctx, authClaims, id)
	})
}
