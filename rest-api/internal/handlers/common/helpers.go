package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
)

func decodeAndValidateRequestBody[T any](request *http.Request, requestBodyStruct *T) error {
	if err := json.NewDecoder(request.Body).Decode(&requestBodyStruct); err != nil {
		return fmt.Errorf("Bad request body")
	}

	validate := validator.New()
	if err := validate.Struct(requestBodyStruct); err != nil {
		return fmt.Errorf("Invalid parameters. Check documentation for further information.")
	}
	return nil
}

func getIdFromPath(request *http.Request) (int, error) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("Invalid ID", http.StatusBadRequest)
	}
	return id, nil
}
