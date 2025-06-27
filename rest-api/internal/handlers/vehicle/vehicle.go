package vehicle

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/czxrny/veh-sense-backend/shared/database"
	"github.com/czxrny/veh-sense-backend/shared/models"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
)

func GetVehicles(response http.ResponseWriter, request *http.Request) {
	db := database.GetDatabaseClient()
	/* TODO - RETURN ONLY THE VEHICLES FROM THE ORGRANIZATION/PRIVATE OWNER */
	/* POSSIBLE USAGE - OWNER WILL HAVE MULTIPLE VEHICLES THAT CAN BE DISPLAYED UPON THE START OF THE APP */
	var vehicles []models.Vehicle

	if err := db.Find(&vehicles).Error; err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(vehicles)
}

func AddVehicle(response http.ResponseWriter, request *http.Request) {
	db := database.GetDatabaseClient()
	var newVehicle models.Vehicle
	if err := json.NewDecoder(request.Body).Decode(&newVehicle); err != nil {
		http.Error(response, "Invalid input", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	validate := validator.New()
	if err := validate.Struct(newVehicle); err != nil {
		http.Error(response, "Invalid parameters. Check documentation for further information.", http.StatusBadRequest)
		return
	}

	if err := db.Create(&newVehicle).Error; err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(newVehicle)
}

func GetVehicleById(response http.ResponseWriter, request *http.Request) {
	/* PROVIDE ONLY IF THE USER IS THE OWNER! */
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(response, "Invalid ID", http.StatusBadRequest)
		return
	}
	db := database.GetDatabaseClient()
	var vehicle models.Vehicle

	if err := db.First(&vehicle, id).Error; err != nil {
		http.Error(response, "Not found", http.StatusNotFound)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(vehicle)
}

func UpdateVehicle(response http.ResponseWriter, request *http.Request) {
	/* TODO - ONLY THE ORGANIZATION ADMIN / OWNER CAN EDIT THE VEHICLE INFO.. */
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(response, "Invalid ID", http.StatusBadRequest)
		return
	}

	db := database.GetDatabaseClient()

	var VehicleUpdate models.VehicleUpdate
	if err := json.NewDecoder(request.Body).Decode(&VehicleUpdate); err != nil {
		http.Error(response, "Invalid input", http.StatusBadRequest)
		return
	}

	result := db.Model(&models.Vehicle{}).Where("id=?", id).Updates(VehicleUpdate)
	if result.Error != nil {
		http.Error(response, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(response, "No rows updated", http.StatusNotFound)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

func DeleteVehicle(response http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(response, "Invalid ID", http.StatusBadRequest)
		return
	}

	db := database.GetDatabaseClient()
	result := db.Delete(&models.Vehicle{}, id)
	if result.Error != nil {
		http.Error(response, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(response, "No record found to delete", http.StatusNotFound)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}
