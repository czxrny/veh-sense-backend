package user

import (
	"fmt"

	"github.com/czxrny/veh-sense-backend/shared/auth"
	"github.com/czxrny/veh-sense-backend/shared/database"
	"github.com/czxrny/veh-sense-backend/shared/models"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(userRegisterInfo *models.UserRegisterInfo, organizationId *int, role string) (*models.UserTokenResponse, error) {
	var err error
	db := database.GetDatabaseClient()

	var resultAuth []models.UserAuth
	db.Where("email = ?", userRegisterInfo.Email).Find(&resultAuth)
	if len(resultAuth) > 0 {
		return nil, fmt.Errorf("Email is already taken.")
	}

	userRegisterInfo.Password, err = auth.EncryptThePassword(userRegisterInfo.Password)
	if err != nil {
		return nil, err
	}

	newUser := models.UserAuth{
		Email:    userRegisterInfo.Email,
		Password: userRegisterInfo.Password,
		Role:     role,
	}

	if err := db.Create(&newUser).Error; err != nil {
		return nil, err
	}

	userInfo := models.UserInfo{
		ID:              newUser.ID,
		UserName:        userRegisterInfo.UserName,
		OrganizationId:  organizationId,
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
}

func LoginUser(userCredentials *models.UserCredentials) (*models.UserTokenResponse, error) {
	db := database.GetDatabaseClient()

	var userAuth models.UserAuth
	db.Where("email = ?", userCredentials.Email).Find(&userAuth)
	if userAuth.ID == 0 {
		return nil, fmt.Errorf("User does not exist.")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userCredentials.Password), []byte(userAuth.Password)); err != nil {
		return nil, fmt.Errorf("Invalid login credentials.")
	}

	var userInfo models.UserInfo
	db.Where("id = ?", userAuth.ID).Find(&userInfo)

	token, err := auth.CreateToken(userAuth, userInfo)
	if err != nil {
		return nil, err
	}

	return &models.UserTokenResponse{
		Token:   token,
		LocalId: userAuth.ID,
	}, nil
}

func UpdateLoginCredentials(credUpdateRequest *models.UserCredentialsUpdateRequest) (*models.UserTokenResponse, error) {
	db := database.GetDatabaseClient()

	var userAuth models.UserAuth
	db.Where("email = ?", credUpdateRequest.Email).Find(&userAuth)
	if userAuth.ID == 0 {
		return nil, fmt.Errorf("User does not exist.")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(credUpdateRequest.Password), []byte(userAuth.Password)); err != nil {
		return nil, fmt.Errorf("Invalid login credentials.")
	}

	var tempAuth models.UserAuth
	db.Where("email = ?", credUpdateRequest.NewEmail).Find(&tempAuth)
	if userAuth.ID != 0 {
		return nil, fmt.Errorf("User email is already taken.")
	}

	if credUpdateRequest.NewEmail == "" || credUpdateRequest.NewPassword == "" {
		return nil, fmt.Errorf("User should pass email or password, or both to update")
	}

	if credUpdateRequest.NewEmail != "" {
		userAuth.Email = credUpdateRequest.NewEmail
	}
	if credUpdateRequest.NewPassword != "" {
		var err error
		userAuth.Password, err = auth.EncryptThePassword(credUpdateRequest.NewPassword)
		if err != nil {
			return nil, err
		}
	}

	if err := db.Model(&models.UserAuth{}).Where("id=?", userAuth.ID).Updates(userAuth).Error; err != nil {
		return nil, err
	}

	var userInfo models.UserInfo
	db.Where("id = ?", userAuth.ID).Find(&userInfo)

	token, err := auth.CreateToken(userAuth, userInfo)
	if err != nil {
		return nil, err
	}

	return &models.UserTokenResponse{
		Token:   token,
		LocalId: userAuth.ID,
	}, nil
}
