package handler

import (
	"context"
	"fmt"
	"net/http"

	s "github.com/czxrny/veh-sense-backend/batch-receiver/internal/domain/upload/service"
	"github.com/czxrny/veh-sense-backend/batch-receiver/internal/model"
	common "github.com/czxrny/veh-sense-backend/shared/handler"
	"github.com/czxrny/veh-sense-backend/shared/middleware"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

type UploadHandler struct {
	s *s.Service
}

func NewUploadHandler(service *s.Service) *UploadHandler {
	return &UploadHandler{
		s: service,
	}
}

func (uh *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	common.PostHandlerSilent(w, r, func(ctx context.Context, request *model.UploadRideRequest) error {
		authClaims, ok := ctx.Value(middleware.AuthKeyName).(models.AuthInfo)
		if !ok {
			return fmt.Errorf("Error: Internal server error. Something went wrong while decoding the JWT.")
		}

		_, err := uh.s.UploadRide(ctx, authClaims, *request)
		return err
	})
}
