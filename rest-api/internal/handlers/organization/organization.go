package organization

import (
	"context"
	"fmt"
	"net/http"

	"github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/common"
	"github.com/czxrny/veh-sense-backend/shared/database"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

// Root only
func GetAllOrganizations(w http.ResponseWriter, r *http.Request) {
	common.GetAllHandler(w, r, func(ctx context.Context) ([]models.Organization, error) {
		authClaims, ok := ctx.Value("authClaims").(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		if authClaims.Role != "root" {
			return nil, fmt.Errorf("Error: User is unauthorized to view the assets.")
		}

		db := database.GetDatabaseClient()

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
		authClaims, ok := ctx.Value("authClaims").(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		if authClaims.Role != "root" {
			return nil, fmt.Errorf("Error: User is unauthorized to add the asset.")
		}

		db := database.GetDatabaseClient()

		if err := db.Create(organization).Error; err != nil {
			return nil, err
		}

		return organization, nil
	})
}

// Root only
func DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	common.DeleteHandler(w, r, func(ctx context.Context, id int) error {
		authClaims, ok := ctx.Value("authClaims").(models.AuthInfo)
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
