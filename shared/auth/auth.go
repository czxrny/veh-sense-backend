package auth

import (
	"fmt"

	"github.com/czxrny/veh-sense-backend/shared/models"
	"golang.org/x/crypto/bcrypt"
)

func EncryptThePassword(userRegisterInfo *models.UserRegisterInfo) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(userRegisterInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Error hashing password")
	}
	userRegisterInfo.Password = string(hashed)
	return nil
}

func ValidatePasswordAndReturnToken(userInfo models.UserAuth, password string) (string, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("Invalid login credentials.")
	}

	token, err := CreateToken(userInfo)
	if err != nil {
		return "", err
	}

	return token, nil
}

// UserAuth is able to edit asset if he was the one posting, or has admin role.
func IsAuthorizedToEditAsset(token string, originalUser string) error {
	UserAuth, err := ExtractFromToken(token, "usr")
	if err != nil {
		return err
	}

	role, err := ExtractFromToken(token, "rol")
	if err != nil {
		return err
	}

	if UserAuth != originalUser && role != "admin" {
		return fmt.Errorf("UserAuth is unauthorized.")
	}

	return nil
}
