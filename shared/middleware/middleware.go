package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/czxrny/veh-sense-backend/shared/auth"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

type ContextKey string

const AuthKeyName ContextKey = "authClaims"

func RequireAPIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		fmt.Println(authHeader)
		if !strings.HasPrefix(authHeader, "ApiKey ") {
			http.Error(w, fmt.Errorf("Missing or malformed Authorization header").Error(), http.StatusUnauthorized)
			return
		}

		apiKey := strings.TrimPrefix(authHeader, "ApiKey ")

		if apiKey != os.Getenv("API_KEY") {
			http.Error(w, fmt.Errorf("Invalid Key in Authorization Header").Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func JWTClaimsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := readHeaderAndExtractClaims(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		lidFloat, ok := claims["lid"].(float64)
		if !ok {
			http.Error(w, "Invalid token: lid is missing or not a number", http.StatusUnauthorized)
			return
		}

		role, ok := claims["rol"].(string)
		if !ok {
			http.Error(w, "Invalid token: rol is missing or not a string", http.StatusUnauthorized)
			return
		}

		var orgID *int
		orgIDRaw, ok := claims["org"]
		if ok {
			orgFloat, ok := orgIDRaw.(float64)
			if !ok {
				http.Error(w, "Invalid token: org is not a number", http.StatusUnauthorized)
				return
			}
			temp := int(orgFloat)
			orgID = &temp
		}

		authClaims := models.AuthInfo{
			UserID:         int(lidFloat),
			Role:           role,
			OrganizationID: orgID,
		}

		ctx := context.WithValue(r.Context(), AuthKeyName, authClaims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func readHeaderAndExtractClaims(r *http.Request) (map[string]interface{}, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, fmt.Errorf("Missing or malformed Authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if err := auth.VerifyToken(token); err != nil {
		return nil, fmt.Errorf("Bad token! " + err.Error())
	}

	claims, err := auth.ExtractClaimsFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("Invalid token: " + err.Error())
	}

	return claims, nil
}
