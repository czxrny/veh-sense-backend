package user

import (
	"fmt"
	"net/http"

	"github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/common"
	"github.com/czxrny/veh-sense-backend/shared/database"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

func DeleteUserById(response http.ResponseWriter, request *http.Request) {
	// 1:1 like the vehicle implementation ! needs refactor
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
