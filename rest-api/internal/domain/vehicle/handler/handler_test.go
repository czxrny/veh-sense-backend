package handler

import (
	"context"
	"net/http"
	"testing"

	testutils "github.com/czxrny/veh-sense-backend/rest-api/internal/testutils/handlers"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

type mockVehicleService struct {
	findVehiclesFunc func(ctx context.Context, filter models.VehicleFilter) ([]models.Vehicle, error)
	addVehicleFunc   func(ctx context.Context, authInfo models.AuthInfo, vehicle *models.Vehicle) (*models.Vehicle, error)
	getByIdFunc      func(ctx context.Context, authInfo models.AuthInfo, id int) (*models.Vehicle, error)
	updateByIdFunc   func(ctx context.Context, authInfo models.AuthInfo, updatedVehicle *models.VehicleUpdate, id int) (*models.Vehicle, error)
	deleteByIdFunc   func(ctx context.Context, authInfo models.AuthInfo, id int) error
}

func (m *mockVehicleService) FindVehicles(ctx context.Context, filter models.VehicleFilter) ([]models.Vehicle, error) {
	return m.findVehiclesFunc(ctx, filter)
}

func (m *mockVehicleService) AddVehicle(ctx context.Context, authInfo models.AuthInfo, vehicle *models.Vehicle) (*models.Vehicle, error) {
	return m.addVehicleFunc(ctx, authInfo, vehicle)
}

func (m *mockVehicleService) GetById(ctx context.Context, authInfo models.AuthInfo, id int) (*models.Vehicle, error) {
	return m.getByIdFunc(ctx, authInfo, id)
}

func (m *mockVehicleService) UpdateById(ctx context.Context, authInfo models.AuthInfo, updatedVehicle *models.VehicleUpdate, id int) (*models.Vehicle, error) {
	return m.updateByIdFunc(ctx, authInfo, updatedVehicle, id)
}

func (m *mockVehicleService) DeleteById(ctx context.Context, authInfo models.AuthInfo, id int) error {
	return m.deleteByIdFunc(ctx, authInfo, id)
}

const baseUrl = "/vehicles"

func TestGetVehiclesHandler(t *testing.T) {
	mockSvc := &mockVehicleService{
		findVehiclesFunc: func(ctx context.Context, filter models.VehicleFilter) ([]models.Vehicle, error) {
			return []models.Vehicle{
				{ID: 1, Brand: "BMW", Model: "X5"},
				{ID: 2, Brand: "Audi", Model: "A4"},
			}, nil
		},
	}
	handler := NewVehicleHandler(mockSvc)

	testutils.RunBasicGetAllTests(baseUrl, handler.GetVehicles, t)
	customTests := []testutils.HttpTestStruct{
		{"ok query with all params", http.MethodGet, "", baseUrl + "?brand=BMW&model=X5&minEngineCapacity=2000&maxEngineCapacity=3000&minEnginePower=150&maxEnginePower=400&plates=XYZ123", testutils.GetAdminAuth(), nil, http.StatusOK},
		{"bad query params type", http.MethodGet, "", baseUrl + "?minEngineCapacity=abcd", testutils.GetAdminAuth(), nil, http.StatusBadRequest},
	}
	testutils.RunTestCases(customTests, handler.GetVehicles, t)
}

func TestAddVehicleHandler(t *testing.T) {
	const vehicleRequestBodyOK = `{
    "brand": "BMW",
    "model": "X5",
    "year": 2022,
    "engine_capacity": 3000,
    "engine_power": 400,
    "plates": "XYZ123",
    "expected_fuel": 10.5
}`
	const vehicleRequestBodyMissing = `{
    "model": "X5",
    "year": 2022,
    "engine_capacity": 3000,
    "engine_power": 400,
    "plates": "XYZ123",
    "expected_fuel": 10.5
}`

	mockSvc := &mockVehicleService{
		addVehicleFunc: func(ctx context.Context, authInfo models.AuthInfo, vehicle *models.Vehicle) (*models.Vehicle, error) {
			return &models.Vehicle{ID: 1, Brand: "BMW", Model: "X5"}, nil
		},
	}

	handler := NewVehicleHandler(mockSvc)
	testutils.RunBasicAddTests(baseUrl, vehicleRequestBodyOK, vehicleRequestBodyMissing, handler.AddVehicle, t)
}

func TestUpdateVehicleHandler(t *testing.T) {
	const vehicleUpdateRequestBodyOK = `{
    "engine_power": 400,
    "plates": "XYZ123",
    "expected_fuel": 10.5
}`
	const vehicleUpdateRequestBodyBad = `{
    "engine_power": "fvasfaf",
    "plates": "XYZ123",
    "expected_fuel": 10.5
}`
	mockSvc := &mockVehicleService{
		updateByIdFunc: func(ctx context.Context, authInfo models.AuthInfo, updatedVehicle *models.VehicleUpdate, id int) (*models.Vehicle, error) {
			return &models.Vehicle{ID: 1, Brand: "BMW", Model: "X5"}, nil
		},
	}

	handler := NewVehicleHandler(mockSvc)
	testutils.RunBasicUpdateTests(baseUrl, vehicleUpdateRequestBodyOK, vehicleUpdateRequestBodyBad, handler.UpdateVehicle, t)
}

func TestGetVehicleByIdHandler(t *testing.T) {
	mockSvc := &mockVehicleService{
		getByIdFunc: func(ctx context.Context, authInfo models.AuthInfo, id int) (*models.Vehicle, error) {
			return &models.Vehicle{ID: 1, Brand: "BMW", Model: "X5"}, nil
		},
	}

	handler := NewVehicleHandler(mockSvc)
	testutils.RunBasicGetByIdTests(baseUrl, handler.GetVehicleById, t)
}

func TestDeleteVehicleHandler(t *testing.T) {
	mockSvc := &mockVehicleService{
		deleteByIdFunc: func(ctx context.Context, authInfo models.AuthInfo, id int) error {
			return nil
		},
	}
	handler := NewVehicleHandler(mockSvc)
	testutils.RunBasicDeleteByIdTests(baseUrl, handler.DeleteVehicle, t)
}
