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
		claims, err := verifyTokenAndExtractClaims(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		authClaims, err := retrieveAuthClaims(claims)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), AuthKeyName, authClaims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func retrieveAuthClaims(claims map[string]interface{}) (models.AuthInfo, error) {
	lidFloat, ok := claims["lid"].(float64)
	if !ok {
		return models.AuthInfo{}, fmt.Errorf("Invalid token: lid is missing or not a number")
	}

	role, ok := claims["rol"].(string)
	if !ok {
		return models.AuthInfo{}, fmt.Errorf("Invalid token: rol is missing or not a string")
	}

	var orgID *int
	orgIDRaw, ok := claims["org"]
	if ok {
		orgFloat, ok := orgIDRaw.(float64)
		if !ok {
			return models.AuthInfo{}, fmt.Errorf("Invalid token: org is not a number")
		}
		temp := int(orgFloat)
		orgID = &temp
	}

	return models.AuthInfo{
		UserID:         int(lidFloat),
		Role:           role,
		OrganizationID: orgID,
	}, nil
}

func verifyTokenAndExtractClaims(r *http.Request) (map[string]interface{}, error) {
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
