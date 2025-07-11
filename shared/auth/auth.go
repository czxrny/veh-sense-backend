package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func EncryptThePassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Error hashing password")
	}
	return string(hashed), nil
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
