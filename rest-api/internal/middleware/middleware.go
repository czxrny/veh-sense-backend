package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/czxrny/veh-sense-backend/shared/auth"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

type contextKey string

const authKeyName contextKey = "authClaims"

func JWTClaimsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(response, "Missing or malformed Authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if err := auth.VerifyToken(token); err != nil {
			http.Error(response, "Bad token! "+err.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := auth.ExtractClaimsFromToken(token)
		if err != nil {
			http.Error(response, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		lidFloat, ok := claims["lid"].(float64)
		if !ok {
			http.Error(response, "Invalid token: lid is missing or not a number", http.StatusUnauthorized)
			return
		}

		role, ok := claims["rol"].(string)
		if !ok {
			http.Error(response, "Invalid token: rol is missing or not a string", http.StatusUnauthorized)
			return
		}

		orgFloat, ok := claims["org"].(float64)
		if !ok {
			http.Error(response, "Invalid token: org is missing or not a number", http.StatusUnauthorized)
			return
		}

		authClaims := models.AuthInfo{
			UserID:         int(lidFloat),
			Role:           role,
			OrganizationID: int(orgFloat),
		}

		ctx := context.WithValue(request.Context(), authKeyName, authClaims)
		next.ServeHTTP(response, request.WithContext(ctx))
	})
}
