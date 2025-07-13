package raport

import (
	"context"
	"fmt"
	"net/http"

	"github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/common"
	"github.com/czxrny/veh-sense-backend/shared/database"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

func GetRaport(w http.ResponseWriter, r *http.Request) {
	common.GetAllHandler(w, r, func(ctx context.Context) ([]models.Raport, error) {
		authClaims, ok := ctx.Value("authClaims").(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		db := database.GetDatabaseClient()

		switch authClaims.Role {
		case "user":
			db = db.Where("user_id = ?", authClaims.UserID)
		case "admin":
			db = db.Where("ogranization_id = ?", authClaims.OrganizationID)
		}

		var raports []models.Raport
		if err := db.Find(&raports).Error; err != nil {
			return nil, err
		}

		return raports, nil
	})
}
