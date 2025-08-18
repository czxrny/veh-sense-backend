package vehicle

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/common"
	"github.com/czxrny/veh-sense-backend/rest-api/internal/middleware"
	"github.com/czxrny/veh-sense-backend/shared/database"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

// Returns all the vehicles that are available to user:
// If the user is a private user: returns all of the private cars
// If the user is corporate and is not an admin: returns all of the assigned cars + shared (the ones without the owner_id)
// If the user is an admin of organization: returns all of the organization cars
// Root gets all the info.
func GetVehicles(w http.ResponseWriter, r *http.Request) {
	common.GetAllHandler(w, r, func(ctx context.Context, query url.Values) ([]models.Vehicle, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		db := database.GetDatabaseClient()

		if query.Has("brand") {
			db = db.Where("brand = ?", query.Get("brand"))
		}
		if query.Has("minCapacity") {
			db = db.Where("engine_capacity >= ?", query.Get("minEngineCapacity"))
		}
		if query.Has("maxCapacity") {
			db = db.Where("engine_capacity <= ?", query.Get("maxEngineCapacity"))
		}
		if query.Has("minEnginePower") {
			db = db.Where("engine_power >= ?", query.Get("minEnginePower"))
		}
		if query.Has("maxEnginePower") {
			db = db.Where("engine_power <= ?", query.Get("maxEnginePower"))
		}

		switch authClaims.Role {
		case "user":
			if authClaims.OrganizationID != nil {
				db = db.Where("(owner_id = ? OR (organization_id = ? AND owner_id IS NULL))",
					authClaims.UserID,
					authClaims.OrganizationID,
				)
			} else {
				db = db.Where("owner_id = ?", authClaims.UserID)
			}
		case "admin":
			db = db.Where("organization_id = ?", authClaims.OrganizationID)
		}

		var vehicles []models.Vehicle
		if err := db.Find(&vehicles).Error; err != nil {
			return nil, err
		}

		return vehicles, nil
	})
}

func AddVehicle(w http.ResponseWriter, r *http.Request) {
	common.PostHandler(w, r, func(ctx context.Context, vehicle *models.Vehicle) (*models.Vehicle, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		switch authClaims.Role {
		case "user":
			if authClaims.OrganizationID != nil {
				return nil, fmt.Errorf("Error: Corporate user is unauthorized to add new vehicles to fleet. Login as an admin to proceed.")
			}
			vehicle.OwnerID = &authClaims.UserID
			vehicle.OrganizationID = authClaims.OrganizationID

		// Admin sets the vehicle organization_id automatically, and can pass the owner_id to specify the user that will be the owner
		case "admin":
			vehicle.OrganizationID = authClaims.OrganizationID

		case "root":
			if vehicle.OwnerID == nil || vehicle.OrganizationID == nil {
				return nil, fmt.Errorf("Error: Bad Request: Please specify either the organization_id or owner_id to proceed")
			}
		}

		vehicle.ID = 0
		db := database.GetDatabaseClient()
		if err := db.Create(vehicle).Error; err != nil {
			return nil, err
		}

		return vehicle, nil
	})
}

// Returns the vehicle - only if the user is the owner of the vehicle / the vehicle is shared in the corporation the user is in
func GetVehicleById(w http.ResponseWriter, r *http.Request) {
	common.GetByIdHandler(w, r, func(ctx context.Context, id int) (*models.Vehicle, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		db := database.GetDatabaseClient()

		var vehicle models.Vehicle
		if err := db.First(&vehicle, id).Error; err != nil {
			return nil, err
		}

		isOwner := vehicle.OwnerID != nil && *vehicle.OwnerID == authClaims.UserID
		isShared := vehicle.OrganizationID != nil && authClaims.OrganizationID != nil && *vehicle.OrganizationID == *authClaims.OrganizationID && vehicle.OwnerID == nil
		isOrgAdmin := vehicle.OrganizationID != nil && authClaims.OrganizationID != nil && *vehicle.OrganizationID == *authClaims.OrganizationID && authClaims.Role == "admin"

		if !isOwner && !isShared && !isOrgAdmin && authClaims.Role != "root" {
			return nil, fmt.Errorf("Error: User is unauthorized to view the vehicle.")
		}

		return &vehicle, nil
	})
}

func UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	common.PatchHandler(w, r, func(ctx context.Context, updatedVehicle *models.VehicleUpdate, id int) (*models.Vehicle, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		db := database.GetDatabaseClient()
		var vehicle models.Vehicle
		if err := db.First(&vehicle, id).Error; err != nil {
			return nil, err
		}

		isPrivateOwner := authClaims.OrganizationID == nil && vehicle.OwnerID != nil && *vehicle.OwnerID == authClaims.UserID
		isOrgAdmin := vehicle.OrganizationID != nil && authClaims.OrganizationID != nil && *vehicle.OrganizationID == *authClaims.OrganizationID && authClaims.Role == "admin"

		if !isPrivateOwner && !isOrgAdmin && authClaims.Role != "root" {
			return nil, fmt.Errorf("Error: User is unauthorized to edit the vehicle.")
		}

		result := db.Model(&models.Vehicle{}).Where("id=?", id).Updates(updatedVehicle)
		if result.Error != nil {
			return nil, result.Error
		}

		if err := db.First(&vehicle, id).Error; err != nil {
			return nil, err
		}

		return &vehicle, nil
	})
}

func DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	common.DeleteHandler(w, r, func(ctx context.Context, id int) error {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		db := database.GetDatabaseClient()
		var vehicle models.Vehicle
		if err := db.First(&vehicle, id).Error; err != nil {
			return err
		}

		isPrivateOwner := authClaims.OrganizationID == nil && vehicle.OwnerID != nil && *vehicle.OwnerID == authClaims.UserID
		isOrgAdmin := vehicle.OrganizationID != nil && authClaims.OrganizationID != nil && *vehicle.OrganizationID == *authClaims.OrganizationID && authClaims.Role == "admin"

		if !isPrivateOwner && !isOrgAdmin && authClaims.Role != "root" {
			return fmt.Errorf("Error: User is unauthorized to delete the vehicle.")
		}

		result := db.Delete(&models.Vehicle{}, id)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
}
