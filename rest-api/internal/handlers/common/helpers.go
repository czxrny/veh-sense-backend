package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
)

func decodeAndValidateRequestBody[T any](r *http.Request, requestBodyStruct *T) error {
	if err := json.NewDecoder(r.Body).Decode(&requestBodyStruct); err != nil {
		return fmt.Errorf("Bad r body")
	}

	validate := validator.New()
	if err := validate.Struct(requestBodyStruct); err != nil {
		return fmt.Errorf("Invalid parameters. Check documentation for further information.")
	}
	return nil
}

func getIdFromPath(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("Invalid ID", http.StatusBadRequest)
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
