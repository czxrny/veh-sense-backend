package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/common"
	userService "github.com/czxrny/veh-sense-backend/rest-api/internal/services/user"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

func RegisterPrivateUser(w http.ResponseWriter, r *http.Request) {
	common.PostHandler(w, r, func(ctx context.Context, userRegisterInfo *models.UserRegisterInfo) (*models.UserTokenResponse, error) {
		return userService.RegisterUser(userRegisterInfo, nil, "user")
	})
}

func RegisterCorporateUser(w http.ResponseWriter, r *http.Request) {
	common.PostHandler(w, r, func(ctx context.Context, userRegisterInfo *models.UserRegisterInfo) (*models.UserTokenResponse, error) {
		authClaims, ok := ctx.Value("authClaims").(models.AuthInfo)
		if !ok || authClaims.Role != "admin" {
			return nil, fmt.Errorf("Error: to create an organization user, login as an admin and pass the JWT!")
		}
		return userService.RegisterUser(userRegisterInfo, authClaims.OrganizationID, authClaims.Role)
	})
}

func RegisterUserRoot(w http.ResponseWriter, r *http.Request) {
	common.PostHandler(w, r, func(ctx context.Context, userRegisterInfo *models.UserRegisterInfoRoot) (*models.UserTokenResponse, error) {
		authClaims, ok := ctx.Value("authClaims").(models.AuthInfo)
		if !ok || authClaims.Role != "root" {
			return nil, fmt.Errorf("Error: to create a custom user, login as a root and pass the JWT!")
		}

		userInfo := models.UserRegisterInfo{
			UserName: userRegisterInfo.UserName,
			Email:    userRegisterInfo.Email,
			Password: userRegisterInfo.Password,
		}

		return userService.RegisterUser(&userInfo, userRegisterInfo.OrganizationID, userRegisterInfo.Role)
	})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	common.PostHandler(w, r, func(ctx context.Context, userCredentials *models.UserCredentials) (*models.UserTokenResponse, error) {
		return userService.LoginUser(userCredentials)
	})
}

// Requires the user to login and pass updated information
func UpdateLoginCredentials(w http.ResponseWriter, r *http.Request) {
	common.PostHandler(w, r, func(ctx context.Context, credUpdateRequest *models.UserCredentialsUpdateRequest) (*models.UserTokenResponse, error) {
		return userService.UpdateLoginCredentials(credUpdateRequest)
	})
}
