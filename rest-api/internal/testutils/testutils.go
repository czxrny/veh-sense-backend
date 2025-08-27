package testutils

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/czxrny/veh-sense-backend/rest-api/internal/middleware"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

func CreateNewRequest(method, url string, authInfo *models.AuthInfo, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, url, body)
	if authInfo != nil {
		ctx := context.WithValue(req.Context(), middleware.AuthKeyName, *authInfo)
		return req.WithContext(ctx)
	}
	return req
}

func GetAdminAuth() *models.AuthInfo {
	return &models.AuthInfo{UserID: 1, OrganizationID: getIntPtr(1), Role: "admin"}
}

func GetUserAuthCorporate() *models.AuthInfo {
	return &models.AuthInfo{UserID: 1, OrganizationID: getIntPtr(1), Role: "user"}
}

func GetUserAuthPrivate() *models.AuthInfo {
	return &models.AuthInfo{UserID: 1, OrganizationID: nil, Role: "user"}
}

func GetRootAuth() *models.AuthInfo {
	return &models.AuthInfo{UserID: 1, OrganizationID: nil, Role: "root"}
}

func getIntPtr(val int) *int { return &val }
