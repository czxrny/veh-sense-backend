package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/common"
	userService "github.com/czxrny/veh-sense-backend/rest-api/internal/services/user"
	"github.com/czxrny/veh-sense-backend/shared/auth"
	"github.com/czxrny/veh-sense-backend/shared/database"
	"github.com/czxrny/veh-sense-backend/shared/models"
	"golang.org/x/crypto/bcrypt"
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
	common.PostHandler(w, r, func(ctx context.Context, userRegisterInfo *models.UserRegisterInfo) (*models.UserTokenResponse, error) {
		db := database.GetDatabaseClient()

		var userAuth models.UserAuth
		db.Where("email = ?", userRegisterInfo.Email).Find(&userAuth)
		if userAuth.ID == 0 {
			return nil, fmt.Errorf("User does not exist.")
		}

		if err := bcrypt.CompareHashAndPassword([]byte(userRegisterInfo.Password), []byte(userAuth.Password)); err != nil {
			return nil, fmt.Errorf("Invalid login credentials.")
		}

		var userInfo models.UserInfo
		db.Where("id = ?", userRegisterInfo.Email).Find(&userInfo)
		if userInfo.ID == 0 {
			return nil, fmt.Errorf("User does not exist.")
		}

		token, err := auth.CreateToken(userAuth, userInfo)
		if err != nil {
			return nil, err
		}

		return &models.UserTokenResponse{
			Token:   token,
			LocalId: userAuth.ID,
		}, nil
	})
}
