package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GenerateRefreshToken() (string, error) {
	length := 32
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

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
