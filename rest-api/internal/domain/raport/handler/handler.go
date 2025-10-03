package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	common "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/common/handler"
	s "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/raport/service"
	"github.com/czxrny/veh-sense-backend/shared/middleware"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

type RaportHandler struct {
	*s.RaportService
}

func NewRaportHandler(raportService *s.RaportService) *RaportHandler {
	return &RaportHandler{
		RaportService: raportService,
	}
}

func (rh *RaportHandler) GetRaports(w http.ResponseWriter, r *http.Request) {
	common.GetAllHandler(w, r, func(ctx context.Context, query url.Values) ([]models.Raport, error) {
		authClaims, ok := ctx.Value("authClaims").(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		filter := &models.RaportFilter{
			CreatedAfter:   query.Get("createdAfter"),
			CreatedBefore:  query.Get("createdBefore"),
			UserID:         authClaims.UserID,
			OrganizationID: authClaims.OrganizationID,
			Role:           authClaims.Role,
		}

		return rh.RaportService.FindAllReports(ctx, *filter)
	})
}

func (rh *RaportHandler) DeleteRaport(w http.ResponseWriter, r *http.Request) {
	common.DeleteHandler(w, r, func(ctx context.Context, id int) error {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		return rh.RaportService.DeleteById(ctx, authClaims, id)
	})
}
