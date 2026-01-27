package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	s "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/report/service"
	common "github.com/czxrny/veh-sense-backend/shared/handler"
	"github.com/czxrny/veh-sense-backend/shared/middleware"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

type ReportHandler struct {
	*s.ReportService
}

func NewReportHandler(ReportService *s.ReportService) *ReportHandler {
	return &ReportHandler{
		ReportService: ReportService,
	}
}

func (rh *ReportHandler) GetReports(w http.ResponseWriter, r *http.Request) {
	common.GetAllHandler(w, r, func(ctx context.Context, query url.Values) ([]models.Report, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		filter := &models.ReportFilter{
			CreatedAfter:   query.Get("createdAfter"),
			CreatedBefore:  query.Get("createdBefore"),
			UserID:         authClaims.UserID,
			OrganizationID: authClaims.OrganizationID,
			Role:           authClaims.Role,
		}

		return rh.ReportService.FindAllReports(ctx, *filter)
	})
}

func (rh *ReportHandler) GetReportDataById(w http.ResponseWriter, r *http.Request) {
	common.GetByIdHandler(w, r, func(ctx context.Context, id int) (*models.RideRecord, error) {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return nil, fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		return rh.ReportService.GetRideData(ctx, authClaims, id)
	})
}

func (rh *ReportHandler) DeleteReport(w http.ResponseWriter, r *http.Request) {
	common.DeleteHandler(w, r, func(ctx context.Context, id int) error {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		return rh.ReportService.DeleteById(ctx, authClaims, id)
	})
}
