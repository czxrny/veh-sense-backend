package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/common"
	"github.com/czxrny/veh-sense-backend/shared/auth"
	"github.com/czxrny/veh-sense-backend/shared/database"
	"github.com/czxrny/veh-sense-backend/shared/models"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	common.PostHandler(w, r, func(ctx context.Context, userRegisterInfo *models.UserRegisterInfo) (*models.UserTokenResponse, error) {
		db := database.GetDatabaseClient()

		var resultAuth []models.UserAuth
		db.Where("email = ?", userRegisterInfo.Email).Find(&resultAuth)
		if len(resultAuth) > 0 {
			return nil, fmt.Errorf("Email is already taken.")
		}

		if err := auth.EncryptThePassword(userRegisterInfo); err != nil {
			return nil, err
		}
		newUser := models.UserAuth{
			Email:    userRegisterInfo.Email,
			Password: userRegisterInfo.Password,
			Role:     "user", // by default
		}
		if err := db.Create(&newUser).Error; err != nil {
			return nil, err
		}

		// to do - can only add an user to a organization if you are the organizations admin! maybe through the jwt? not from request body?
		// if the jwt is not passed - then just skip? and if it is passed - check if it is an admin? then check the corporation id?
		userInfo := models.UserInfo{
			ID:              newUser.ID,
			UserName:        userRegisterInfo.UserName,
			OrganizationId:  nil,
			TotalKilometers: 0,
		}
		if err := db.Create(&userInfo).Error; err != nil {
			return nil, err
		}

		token, err := auth.CreateToken(newUser, userInfo)
		if err != nil {
			return nil, err
		}

		return &models.UserTokenResponse{
			Token:   token,
			LocalId: newUser.ID,
		}, nil
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
