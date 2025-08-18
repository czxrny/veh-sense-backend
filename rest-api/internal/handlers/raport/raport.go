package raport

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

func GetRaports(w http.ResponseWriter, r *http.Request) {
	common.GetAllHandler(w, r, func(ctx context.Context, query url.Values) ([]models.Raport, error) {
		authClaims, ok := ctx.Value("authClaims").(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		db := database.GetDatabaseClient()

		if query.Has("createdAfter") {
			db = db.Where("frame_time >= ?", query.Get("createdAfter"))
		}
		if query.Has("createdBefore") {
			db = db.Where("frame_time <= ?", query.Get("createdBefore"))
		}

		switch authClaims.Role {
		case "user":
			db = db.Where("user_id = ?", authClaims.UserID)
		case "admin":
			db = db.Where("ogranization_id = ?", authClaims.OrganizationID)
		}

		var raports []models.Raport
		if err := db.Find(&raports).Error; err != nil {
			return nil, err
		}

		return raports, nil
	})
}

func DeleteRaport(w http.ResponseWriter, r *http.Request) {
	common.DeleteHandler(w, r, func(ctx context.Context, id int) error {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		db := database.GetDatabaseClient()
		var raport models.Raport
		if err := db.First(&raport, id).Error; err != nil {
			return err
		}

		isPrivateOwner := authClaims.OrganizationID == nil && raport.UserID == authClaims.UserID
		isOrgAdmin := raport.OrganizationID != nil && authClaims.OrganizationID != nil && raport.OrganizationID == authClaims.OrganizationID && authClaims.Role == "admin"

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
