package vehicle

import (
	"encoding/json"
	"net/http"
	"veh-sense-backend/internal/database"
	"veh-sense-backend/internal/models"

	"github.com/go-playground/validator"
)

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
		http.Error(response, "Invalid parameters. Check documentation for further informations.", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO vehicles (
			owner_id, private, brand, model, year, engine_capacity, engine_power, plates, expected_fuel
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	err := db.QueryRow(query,
		newVehicle.OwnerID, /* TO IMPLEMENT - GET THE INFO FROM JWT TOKEN */
		newVehicle.Private, /* TO IMPLEMENT - GET THE INFO FROM JWT TOKEN */
		newVehicle.Brand,
		newVehicle.Model,
		newVehicle.Year,
		newVehicle.EngineCapacity,
		newVehicle.EnginePower,
		newVehicle.Plates,
		newVehicle.ExpectedFuel,
	).Scan(&newVehicle.ID)

	if err != nil {
		http.Error(response, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(newVehicle)
}
