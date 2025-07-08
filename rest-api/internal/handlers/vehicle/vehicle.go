package vehicle

import (
	"fmt"
	"net/http"

	"github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/common"
	"github.com/czxrny/veh-sense-backend/shared/database"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

func GetVehicles(response http.ResponseWriter, request *http.Request) {
	/* TODO - RETURN ONLY THE VEHICLES FROM THE ORGRANIZATION/PRIVATE OWNER */
	/* POSSIBLE USAGE - OWNER WILL HAVE MULTIPLE VEHICLES THAT CAN BE DISPLAYED UPON THE START OF THE APP */
	common.GetAllHandler(response, request, func(vehicles *[]models.Vehicle) error {
		db := database.GetDatabaseClient()
		return db.Find(&vehicles).Error
	})
}

func AddVehicle(response http.ResponseWriter, request *http.Request) {
	common.PostHandler(response, request, func(response http.ResponseWriter, request *http.Request, vehicle *models.Vehicle) (*models.Vehicle, error) {
		db := database.GetDatabaseClient()
		if err := db.Create(vehicle).Error; err != nil {
			return nil, err
		}

		return vehicle, nil
	})
}

func GetVehicleById(response http.ResponseWriter, request *http.Request) {
	/* PROVIDE ONLY IF THE USER IS THE OWNER! */
	common.GetByIdHandler(response, request, func(vehicle *models.Vehicle, id int) error {
		db := database.GetDatabaseClient()
		return db.First(&vehicle, id).Error
	})
}

func UpdateVehicle(response http.ResponseWriter, request *http.Request) {
	/* TODO - ONLY THE ORGANIZATION ADMIN / OWNER CAN EDIT THE VEHICLE INFO.. */
	common.PatchHandler(response, request, func(vehicle *models.VehicleUpdate, id int) error {
		db := database.GetDatabaseClient()
		result := db.Model(&models.Vehicle{}).Where("id=?", id).Updates(vehicle)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("No rows updated")
		}

		return nil
	})
}

func DeleteVehicle(response http.ResponseWriter, request *http.Request) {
	common.DeleteHandler(response, request, func(id int) error {
		db := database.GetDatabaseClient()
		result := db.Delete(&models.Vehicle{}, id)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("No record found to delete")
		}

		return nil
	})
}
