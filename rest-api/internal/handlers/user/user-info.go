package user

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

// For admins
func GetAllUsersInfo(w http.ResponseWriter, r *http.Request) {
	common.GetAllHandler(w, r, func(ctx context.Context, query url.Values) ([]models.UserInfo, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		if authClaims.Role != "admin" {
			return nil, fmt.Errorf("Error: User is unauthorized to see other users.")
		}

		db := database.GetDatabaseClient()

		var users []models.UserInfo
		if err := db.Where("organization_id = ?", authClaims.OrganizationID).Find(&users).Error; err != nil {
			return nil, err
		}

		return users, nil
	})
}

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

// Must be either owner, admin of the user org, or root
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
