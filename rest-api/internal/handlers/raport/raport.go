package raport

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/czxrny/veh-sense-backend/rest-api/internal/handlers/common"
	"github.com/czxrny/veh-sense-backend/shared/database"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

func GetRaports(w http.ResponseWriter, r *http.Request) {
	common.GetAllHandler(w, r, func(ctx context.Context, query url.Values) ([]models.Raport, error) {
		authClaims, ok := ctx.Value("authClaims").(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		db := database.GetDatabaseClient()

		if query.Has("createdAfter") {
			db = db.Where("frame_time >= ?", query.Get("createdAfter"))
		}
		if query.Has("createdBefore") {
			db = db.Where("frame_time <= ?", query.Get("createdBefore"))
		}

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
