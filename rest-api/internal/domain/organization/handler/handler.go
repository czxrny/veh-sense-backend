package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	common "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/common/handler"
	o "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/organization/service"
	"github.com/czxrny/veh-sense-backend/shared/middleware"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

type OrganizationHandler struct {
	*o.OrganizationService
}

func NewOrganizationHandler(organizationService *o.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		OrganizationService: organizationService,
	}
}

// Root only
func (s *OrganizationHandler) GetAllOrganizations(w http.ResponseWriter, r *http.Request) {
	common.GetAllHandler(w, r, func(ctx context.Context, query url.Values) ([]models.Organization, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		filter := models.OrganizationFilter{
			City:    query.Get("city"),
			Country: query.Get("country"),
		}

		return s.OrganizationService.FindOrganizations(ctx, authClaims, filter)
	})
}

// Root only
func (s *OrganizationHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	common.PostHandler(w, r, func(ctx context.Context, organization *models.Organization) (*models.Organization, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		return s.OrganizationService.AddOrganization(ctx, organization, authClaims)
	})
}

// Root only
func (s *OrganizationHandler) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	common.DeleteHandler(w, r, func(ctx context.Context, id int) error {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		return s.OrganizationService.DeleteById(ctx, authClaims, id)
	})
}

// Returns only the info about the organization of the JWT user...
func (s *OrganizationHandler) GetMyOrganizationInfo(w http.ResponseWriter, r *http.Request) {
	common.GetSimpleHandler(w, r, func(ctx context.Context) (*models.Organization, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		return s.OrganizationService.GetMyOrganizationInfo(ctx, authClaims)
	})
}

// Based on JWT Token... Requires admin role of organization / root.
func (s *OrganizationHandler) UpdateMyOrganization(w http.ResponseWriter, r *http.Request) {
	common.PatchSimpleHandler(w, r, func(ctx context.Context, organizationUpdate *models.OrganizationUpdate) (*models.Organization, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		return s.OrganizationService.UpdateMyOrganization(ctx, authClaims, organizationUpdate)
	})
}
