package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/czxrny/veh-sense-backend/shared/middleware"
	"github.com/czxrny/veh-sense-backend/shared/models"
	"github.com/go-chi/chi"
)

type HttpTestStruct struct {
	Name       string
	Method     string
	ID         string
	Url        string
	AuthInfo   *models.AuthInfo
	Body       io.Reader
	WantStatus int
}

func CreateNewRequest(method, url string, authInfo *models.AuthInfo, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, url, body)
	if authInfo != nil {
		ctx := context.WithValue(req.Context(), middleware.AuthKeyName, *authInfo)
		return req.WithContext(ctx)
	}
	return req
}

func AddChiIdToContext(r *http.Request, id string) *http.Request {
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, chi.NewRouteContext())
	chiCtx := chi.RouteContext(ctx)
	chiCtx.URLParams.Add("id", id)

	return r.WithContext(ctx)
}

func GetRootAuth() *models.AuthInfo {
	return &models.AuthInfo{UserID: 1, OrganizationID: nil, Role: "root"}
}

func GetAdminAuth() *models.AuthInfo {
	return &models.AuthInfo{UserID: 2, OrganizationID: getIntPtr(1), Role: "admin"}
}

// from the same organization
func GetUserAuthCorporate() *models.AuthInfo {
	return &models.AuthInfo{UserID: 3, OrganizationID: getIntPtr(1), Role: "user"}
}

func GetUserAuthPrivate() *models.AuthInfo {
	return &models.AuthInfo{UserID: 4, OrganizationID: nil, Role: "user"}
}

func getIntPtr(val int) *int { return &val }
