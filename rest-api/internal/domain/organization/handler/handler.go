package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	database "github.com/czxrny/veh-sense-backend/rest-api/internal/app"
	common "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/common/handler"
	"github.com/czxrny/veh-sense-backend/rest-api/internal/middleware"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

// Root only
func GetAllOrganizations(w http.ResponseWriter, r *http.Request) {
	common.GetAllHandler(w, r, func(ctx context.Context, query url.Values) ([]models.Organization, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		if authClaims.Role != "root" {
			return nil, fmt.Errorf("Error: User is unauthorized to view the assets.")
		}

		db := database.GetDatabaseClient()

		if query.Has("city") {
			db = db.Where("city = ?", query.Get("city"))
		}
		if query.Has("country") {
			db = db.Where("country = ?", query.Get("country"))
		}

		var organizations []models.Organization
		if err := db.Find(&organizations).Error; err != nil {
			return nil, err
		}

		return organizations, nil
	})
}

// Root only
func CreateOrganization(w http.ResponseWriter, r *http.Request) {
	common.PostHandler(w, r, func(ctx context.Context, organization *models.Organization) (*models.Organization, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		if authClaims.Role != "root" {
			return nil, fmt.Errorf("Error: User is unauthorized to add the asset.")
		}

		db := database.GetDatabaseClient()

		organization.ID = 0
		if err := db.Create(organization).Error; err != nil {
			return nil, err
		}

		return organization, nil
	})
}

// Root only
func DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	common.DeleteHandler(w, r, func(ctx context.Context, id int) error {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		if authClaims.Role != "root" {
			return fmt.Errorf("Error: User is unauthorized to delete the asset.")
		}

		db := database.GetDatabaseClient()

		return db.Delete(&models.Organization{}, id).Error
	})
}

// Returns only the info about the organization of the JWT user...
func GetMyOrganizationInfo(w http.ResponseWriter, r *http.Request) {
	common.GetSimpleHandler(w, r, func(ctx context.Context) (*models.Organization, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		if authClaims.OrganizationID == nil {
			return nil, fmt.Errorf("Error: User does not have an organization.")
		}

		db := database.GetDatabaseClient()

		var organization models.Organization
		if err := db.First(&organization, authClaims.OrganizationID).Error; err != nil {
			return nil, err
		}

		return &organization, nil
	})
}

// Based on JWT Token... Requires admin role of organization / root.
func PatchMyOrganization(w http.ResponseWriter, r *http.Request) {
	common.PatchSimpleHandler(w, r, func(ctx context.Context, organizationUpdate *models.Organization) (*models.Organization, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		isOrgAdmin := authClaims.OrganizationID != nil && authClaims.Role == "admin"
		if !isOrgAdmin && authClaims.Role != "root" {
			return nil, fmt.Errorf("Error: User is unauthorized to edit the organization.")
		}

		db := database.GetDatabaseClient()

		organizationUpdate.ID = 0
		if err := db.Model(&models.Organization{}).Where("id=?", authClaims.OrganizationID).Updates(organizationUpdate).Error; err != nil {
			return nil, err
		}

		var updatedOrganization models.Organization
		if err := db.First(&updatedOrganization, authClaims.OrganizationID).Error; err != nil {
			return nil, err
		}

		return &updatedOrganization, nil
	})
}
