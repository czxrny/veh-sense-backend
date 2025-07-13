package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/common"
	"github.com/czxrny/veh-sense-backend/rest-api/internal/middleware"
	"github.com/czxrny/veh-sense-backend/shared/database"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

func GetMyUserInfo(w http.ResponseWriter, r *http.Request) {
	common.GetSimpleHandler(w, r, func(ctx context.Context) (*models.UserInfo, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		db := database.GetDatabaseClient()
		var userInfo models.UserInfo
		if err := db.First(&userInfo, authClaims.UserID).Error; err != nil {
			return nil, err
		}

		return &userInfo, nil
	})
}

func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	common.DeleteHandler(w, r, func(ctx context.Context, id int) error {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		db := database.GetDatabaseClient()

		var userInfo models.UserInfo
		if err := db.First(&userInfo, id).Error; err != nil {
			return err
		}

		isOwner := id == authClaims.UserID
		isOrgAdmin := authClaims.OrganizationID != nil && userInfo.OrganizationId != nil && *userInfo.OrganizationId == *authClaims.OrganizationID && authClaims.Role == "admin"

		if !isOwner && !isOrgAdmin && authClaims.Role != "root" {
			return fmt.Errorf("Error: User is unauthorized to delete the user.")
		}

		if err := db.Delete(&models.UserAuth{}, id).Error; err != nil {
			return err
		}
		return db.Delete(&models.UserInfo{}, id).Error
	})
}
