package vehicle

import (
	"context"
	"fmt"
	"net/http"

	"github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/common"
	"github.com/czxrny/veh-sense-backend/shared/database"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

func GetVehicles(w http.ResponseWriter, r *http.Request) {
	/* TODO - RETURN ONLY THE VEHICLES FROM THE ORGRANIZATION/PRIVATE OWNER */
	/* POSSIBLE USAGE - OWNER WILL HAVE MULTIPLE VEHICLES THAT CAN BE DISPLAYED UPON THE START OF THE APP */
	common.GetAllHandler(w, r, func(ctx context.Context) ([]models.Vehicle, error) {
		db := database.GetDatabaseClient()

		var vehicles []models.Vehicle
		if err := db.Find(&vehicles).Error; err != nil {
			return nil, err
		}

		return vehicles, nil
	})
}

func AddVehicle(w http.ResponseWriter, r *http.Request) {
	common.PostHandler(w, r, func(ctx context.Context, vehicle *models.Vehicle) (*models.Vehicle, error) {
		db := database.GetDatabaseClient()
		if err := db.Create(vehicle).Error; err != nil {
			return nil, err
		}

		return vehicle, nil
	})
}

func GetVehicleById(w http.ResponseWriter, r *http.Request) {
	/* PROVIDE ONLY IF THE USER IS THE OWNER! */
	common.GetByIdHandler(w, r, func(ctx context.Context, id int) (*models.Vehicle, error) {
		db := database.GetDatabaseClient()

		var vehicle models.Vehicle
		if err := db.First(&vehicle, id).Error; err != nil {
			return nil, err
		}

		return &vehicle, nil
	})
}

func UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	/* TODO - ONLY THE ORGANIZATION ADMIN / OWNER CAN EDIT THE VEHICLE INFO.. */
	common.PatchHandler(w, r, func(ctx context.Context, updatedVehicle *models.VehicleUpdate, id int) (*models.Vehicle, error) {
		db := database.GetDatabaseClient()
		result := db.Model(&models.Vehicle{}).Where("id=?", id).Updates(updatedVehicle)
		if result.Error != nil {
			return nil, result.Error
		}

		if result.RowsAffected == 0 {
			return nil, fmt.Errorf("Vehicle does not exist!")
		}

		var vehicle models.Vehicle
		if err := db.First(&vehicle, id).Error; err != nil {
			return nil, err
		}

		return &vehicle, nil
	})
}

func DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	common.DeleteHandler(w, r, func(ctx context.Context, id int) error {
		db := database.GetDatabaseClient()
		result := db.Delete(&models.Vehicle{}, id)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("No vehicle found to delete")
		}

		return nil
	})
}
