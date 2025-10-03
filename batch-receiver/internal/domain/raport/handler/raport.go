package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	common "github.com/czxrny/veh-sense-backend/shared/handler"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

func BatchCatcher(w http.ResponseWriter, r *http.Request) {
	common.PostHandlerSilent(w, r, func(ctx context.Context, frame *models.RaportFrame) error {
		fmt.Printf("Timestamp: %s,\nData: %v", time.Now().String(), frame)
		return nil
	})
}
