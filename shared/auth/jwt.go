package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/czxrny/veh-sense-backend/shared/models"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(userAuth models.UserAuth, userInfo models.UserInfo) (string, error) {
	expirationTime := time.Now().Add(time.Hour)

	claims := jwt.MapClaims{
		"lid": userAuth.ID,
		"rol": userAuth.Role,
		"org": userInfo.OrganizationId,
		"iat": time.Now(),
		"exp": expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("token has a wrong signature algorithm set")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return fmt.Errorf("Bad token!: " + err.Error())
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return fmt.Errorf("Invalid token formatting / expired", err)
	}

	return nil
}

func ExtractFromToken(token string, field string) (string, error) {
	claims, err := ExtractClaimsFromToken(token)
	if err != nil {
		return "", err
	}

	fieldVal, ok := claims[field].(string)
	if !ok {
		return "", fmt.Errorf("%v not found or is not a string", field)
	}
	return fieldVal, nil
}

func ExtractClaimsFromToken(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWT format")
	}

	payloadData, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("error decoding JWT payload: %w", err)
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payloadData, &claims); err != nil {
		return nil, fmt.Errorf("invalid JWT payload: %w", err)
	}

	return claims, nil
}
