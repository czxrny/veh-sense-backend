package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/czxrny/veh-sense-backend/shared/apierrors"
	"github.com/czxrny/veh-sense-backend/shared/models"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
)

func handleErrors(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apierrors.ErrBadJWT):
		errorResponse(w, err.Error(), http.StatusUnauthorized)
	case errors.Is(err, apierrors.ErrBadRequest):
		errorResponse(w, err.Error(), http.StatusBadRequest)
	default:
		errorResponse(w, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
	}
}

func errorResponse(response http.ResponseWriter, message string, statusCode int) {
	errorResponse := models.APIError{
		Code:    statusCode,
		Message: message,
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)

	json.NewEncoder(response).Encode(errorResponse)
}

func decodeAndValidateRequestBody[T any](r *http.Request, requestBodyStruct *T) error {
	if err := json.NewDecoder(r.Body).Decode(&requestBodyStruct); err != nil {
		return fmt.Errorf("bad r body")
	}

	validate := validator.New()
	if err := validate.Struct(requestBodyStruct); err != nil {
		return fmt.Errorf("invalid parameters. Check documentation for further information")
	}
	return nil
}

func getIdFromPath(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid ID")
	}
	return id, nil
}

func requestBodyIsEmpty(r *http.Request) bool {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return false
	}
	if len(bodyBytes) > 0 {
		return false
	}

	return true
}
