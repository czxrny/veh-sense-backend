package vehicle

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"veh-sense-backend/internal/database"
	"veh-sense-backend/internal/models"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
)

func GetVehicles(response http.ResponseWriter, request *http.Request) {
	db := database.GetDatabaseClient()
	/* TODO - RETURN ONLY THE VEHICLES FROM THE ORGRANIZATION/PRIVATE OWNER */
	/* POSSIBLE USAGE - OWNER WILL HAVE MULTIPLE VEHICLES THAT CAN BE DISPLAYED UPON THE START OF THE APP */
	rows, err := db.Query("SELECT * FROM vehicles ORDER BY id")
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var vehicles []models.Vehicle
	for rows.Next() {
		var vehicle models.Vehicle
		if err := rows.Scan(
			&vehicle.ID,
			&vehicle.OwnerID,
			&vehicle.Private,
			&vehicle.Brand,
			&vehicle.Model,
			&vehicle.Year,
			&vehicle.EngineCapacity,
			&vehicle.EnginePower,
			&vehicle.Plates,
			&vehicle.ExpectedFuel); err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
		vehicles = append(vehicles, vehicle)
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

func GetVehicleById(response http.ResponseWriter, request *http.Request) {
	/* PROVIDE ONLY IF THE USER IS THE OWNER! */
	id := chi.URLParam(request, "id")
	db := database.GetDatabaseClient()
	var vehicle models.Vehicle

	err := db.QueryRow("SELECT * FROM vehicles WHERE id=$1", id).Scan(
		&vehicle.ID,
		&vehicle.OwnerID,
		&vehicle.Private,
		&vehicle.Brand,
		&vehicle.Model,
		&vehicle.Year,
		&vehicle.EngineCapacity,
		&vehicle.EnginePower,
		&vehicle.Plates,
		&vehicle.ExpectedFuel)
	if err == sql.ErrNoRows {
		http.Error(response, "Item not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(vehicle)
}

func UpdateVehicle(response http.ResponseWriter, request *http.Request) {
	/* TODO - ONLY THE ORGANIZATION ADMIN / OWNER CAN EDIT THE VEHICLE INFO.. */
	id := chi.URLParam(request, "id")
	db := database.GetDatabaseClient()

	var vehicle models.Vehicle
	if err := db.QueryRow("SELECT * FROM vehicles WHERE id=$1", id).Scan(
		&vehicle.ID,
		&vehicle.OwnerID,
		&vehicle.Private,
		&vehicle.Brand,
		&vehicle.Model,
		&vehicle.Year,
		&vehicle.EngineCapacity,
		&vehicle.EnginePower,
		&vehicle.Plates,
		&vehicle.ExpectedFuel); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewDecoder(request.Body).Decode(&vehicle); err != nil {
		http.Error(response, "Invalid input", http.StatusBadRequest)
		return
	}

	query := "UPDATE vehicles SET "
	numberOfParams := 0

	if vehicle.EnginePower >= 0 {
		query += "engine_power=" + strconv.Itoa(vehicle.EnginePower)
		numberOfParams++
	}
	if vehicle.Plates != " " {
		if numberOfParams > 0 {
			query += ", "
		}
		query += "plates=" + vehicle.Plates
	}
	if vehicle.ExpectedFuel >= 0 {
		if numberOfParams > 0 {
			query += ", "
		}
		query += "expected_fuel=" + strconv.Itoa(vehicle.EnginePower)
	}
	query += " WHERE id=" + id

	result, err := db.Exec(query)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(response, "Vehicle not found", http.StatusNotFound)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

func DeleteVehicle(response http.ResponseWriter, request *http.Request) {
	db := database.GetDatabaseClient()
	id := chi.URLParam(request, "id")

	result, err := db.Exec("DELETE FROM vehicles WHERE id=$1", id)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(response, "Item not found", http.StatusNotFound)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}
