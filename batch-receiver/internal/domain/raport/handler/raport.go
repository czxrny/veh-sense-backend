package handler

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/czxrny/veh-sense-backend/batch-receiver/internal/model"
	common "github.com/czxrny/veh-sense-backend/shared/handler"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	common.PostHandlerSilent(w, r, func(ctx context.Context, request *model.UploadRideRequest) error {
		compressed, err := base64.StdEncoding.DecodeString(request.Data)
		if err != nil {
			return fmt.Errorf("invalid base64: %w", err)
		}

		gz, err := gzip.NewReader(bytes.NewReader(compressed))
		if err != nil {
			return fmt.Errorf("invalid gzip: %w", err)
		}
		defer gz.Close()

		var frames []model.ObdFrame
		if err := json.NewDecoder(gz).Decode(&frames); err != nil {
			return fmt.Errorf("invalid OBD JSON: %w", err)
		}

		log.Printf(
			"vehicle_id=%d frames=%d user=%v",
			request.VehicleID,
			len(frames),
			ctx.Value("user_id"),
		)

		return nil
	})
}
