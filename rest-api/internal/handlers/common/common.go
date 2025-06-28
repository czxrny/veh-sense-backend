package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/czxrny/veh-sense-backend/shared/models"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
)

// Parameters:
//   - response: http.ResponseWriter object,
//   - request: *http.Request object,
//   - serviceFunc(*[]T) error - service function writing into the passed generic list all of the found assets.
func GetAllHandler[T any](response http.ResponseWriter, request *http.Request, serviceFunc func(*[]T) error) {
	var assets []T

	if err := serviceFunc(&assets); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(assets)
}

// Parameters:
//   - response: http.ResponseWriter object,
//   - request: *http.Request object,
//   - serviceFunc(*T, int) error - service function that writes into speficied generic structure.
func GetByIdHandler[T any](response http.ResponseWriter, request *http.Request, serviceFunc func(*T, int) error) {
	id, err := getIdFromPath(request)
	if err != nil {
		http.Error(response, err.Error(), http.StatusNotFound)
	}

	var obj T
	if err := serviceFunc(&obj, id); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(obj)
}

// Parameters:
//   - response: http.ResponseWriter object,
//   - request: *http.Request object,
//   - serviceFunc(*T) error - service function puting new data into the database using the generic argument.
func PostHandler[T any](response http.ResponseWriter, request *http.Request, serviceFunc func(*T) error) {
	var obj T
	if err := decodeAndValidateRequestBody(request, &obj); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := serviceFunc(&obj); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(obj)
}

// Parameters:
//   - response: http.ResponseWriter object,
//   - request: *http.Request object,
//   - serviceFunc(*T, int) error - service function updating data of the asset identified by int argument with the data passed as T structure.
func PatchHandler[T any](response http.ResponseWriter, request *http.Request, serviceFunc func(*T, int) error) {
	id, err := getIdFromPath(request)
	if err != nil {
		http.Error(response, err.Error(), http.StatusNotFound)
	}

	var updateObj T
	if err := decodeAndValidateRequestBody(request, &updateObj); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := serviceFunc(&updateObj, id); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

// Parameters:
//   - response: http.ResponseWriter object,
//   - request: *http.Request object,
//   - serviceFunc(int) error - service function deleting an asset identified by the Id from the database.
func DeleteHandler(response http.ResponseWriter, request *http.Request, serviceFunc func(int) error) {
	id, err := getIdFromPath(request)
	if err != nil {
		http.Error(response, err.Error(), http.StatusNotFound)
	}

	if err := serviceFunc(id); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

// Parameters:
//   - response: http.ResponseWriter object,
//   - request: *http.Request object,
//   - ...
func AuthHandler[T any](response http.ResponseWriter, request *http.Request, serviceFunc func(*T) (models.UserTokenResponse, error)) {
	var obj T
	if err := decodeAndValidateRequestBody(request, &obj); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenResponse, err := serviceFunc(&obj)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(tokenResponse)
}

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
