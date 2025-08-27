package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/czxrny/veh-sense-backend/rest-api/internal/testutils"
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

func TestGetVehiclesHandler(t *testing.T) {
	const url = "/vehicles"

	tests := []struct {
		name       string
		url        string
		authInfo   *models.AuthInfo
		body       io.Reader
		wantStatus int
	}{
		{"no token", url, nil, nil, http.StatusUnauthorized},
		{"ok", url, testutils.GetAdminAuth(), nil, http.StatusOK},
		{"with body", url, testutils.GetAdminAuth(), strings.NewReader(`{"invalid":"data"}`), http.StatusBadRequest},
		{"ok query with all params", url + "?brand=BMW&model=X5&minEngineCapacity=2000&maxEngineCapacity=3000&minEnginePower=150&maxEnginePower=400&plates=XYZ123", testutils.GetAdminAuth(), nil, http.StatusOK},
		{"bad query params type", url + "?minEngineCapacity=abcd", testutils.GetAdminAuth(), nil, http.StatusBadRequest},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			mockSvc := &mockVehicleService{
				findVehiclesFunc: func(ctx context.Context, filter models.VehicleFilter) ([]models.Vehicle, error) {
					return []models.Vehicle{
						{ID: 1, Brand: "BMW", Model: "X5"},
						{ID: 2, Brand: "Audi", Model: "A4"},
					}, nil
				},
			}

			req := testutils.CreateNewRequest(http.MethodGet, testCase.url, testCase.authInfo, testCase.body)
			w := httptest.NewRecorder()

			handler := NewVehicleHandler(mockSvc)
			handler.GetVehicles(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != testCase.wantStatus {
				t.Errorf("expected %d, got %d", testCase.wantStatus, res.StatusCode)
			}
		})
	}
}

func TestAddVehicleHandler(t *testing.T) {
	const url = "/vehicles"
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

	tests := []struct {
		name       string
		authInfo   *models.AuthInfo
		body       io.Reader
		wantStatus int
	}{
		{"no token", nil, strings.NewReader(vehicleRequestBodyOK), http.StatusUnauthorized},
		{"ok", testutils.GetAdminAuth(), strings.NewReader(vehicleRequestBodyOK), http.StatusOK},
		{"no request body", testutils.GetAdminAuth(), nil, http.StatusBadRequest},
		{"bad request body", testutils.GetAdminAuth(), strings.NewReader(vehicleRequestBodyMissing), http.StatusBadRequest},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			mockSvc := &mockVehicleService{
				addVehicleFunc: func(ctx context.Context, authInfo models.AuthInfo, vehicle *models.Vehicle) (*models.Vehicle, error) {
					return &models.Vehicle{ID: 1, Brand: "BMW", Model: "X5"}, nil
				},
			}

			req := testutils.CreateNewRequest(http.MethodPost, url, testCase.authInfo, testCase.body)
			w := httptest.NewRecorder()

			handler := NewVehicleHandler(mockSvc)
			handler.AddVehicle(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != testCase.wantStatus {
				t.Errorf("expected %d, got %d", testCase.wantStatus, res.StatusCode)
			}
		})
	}
}

func TestUpdateVehicleHandler(t *testing.T) {
	const url = "/vehicles"
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

	tests := []struct {
		name       string
		id         string
		authInfo   *models.AuthInfo
		body       io.Reader
		wantStatus int
	}{
		{"no token", "1", nil, strings.NewReader(vehicleUpdateRequestBodyOK), http.StatusUnauthorized},
		{"ok", "1", testutils.GetAdminAuth(), strings.NewReader(vehicleUpdateRequestBodyOK), http.StatusOK},
		{"no request body", "1", testutils.GetAdminAuth(), nil, http.StatusBadRequest},
		{"bad request body", "1", testutils.GetAdminAuth(), strings.NewReader(vehicleUpdateRequestBodyBad), http.StatusBadRequest},
		{"bad id in path", "invalid", testutils.GetAdminAuth(), strings.NewReader(vehicleUpdateRequestBodyOK), http.StatusBadRequest},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			mockSvc := &mockVehicleService{
				updateByIdFunc: func(ctx context.Context, authInfo models.AuthInfo, updatedVehicle *models.VehicleUpdate, id int) (*models.Vehicle, error) {
					return &models.Vehicle{ID: 1, Brand: "BMW", Model: "X5"}, nil
				},
			}

			req := testutils.CreateNewRequest(http.MethodPatch, url, testCase.authInfo, testCase.body)
			req = testutils.AddChiIdToContext(req, testCase.id)
			w := httptest.NewRecorder()

			handler := NewVehicleHandler(mockSvc)
			handler.UpdateVehicle(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != testCase.wantStatus {
				t.Errorf("expected %d, got %d", testCase.wantStatus, res.StatusCode)
			}
		})
	}
}

func TestGetVehicleByIdHandler(t *testing.T) {
	const url = "/vehicles"

	tests := []struct {
		name       string
		id         string
		authInfo   *models.AuthInfo
		body       io.Reader
		wantStatus int
	}{
		{"no token", "1", nil, nil, http.StatusUnauthorized},
		{"ok", "1", testutils.GetAdminAuth(), nil, http.StatusOK},
		{"with body", "1", testutils.GetAdminAuth(), strings.NewReader(`{"invalid":"data"}`), http.StatusBadRequest},
		{"bad id in path", "invalid", testutils.GetAdminAuth(), nil, http.StatusBadRequest},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			mockSvc := &mockVehicleService{
				getByIdFunc: func(ctx context.Context, authInfo models.AuthInfo, id int) (*models.Vehicle, error) {
					return &models.Vehicle{ID: 1, Brand: "BMW", Model: "X5"}, nil
				},
			}

			req := testutils.CreateNewRequest(http.MethodPatch, url, testCase.authInfo, testCase.body)
			req = testutils.AddChiIdToContext(req, testCase.id)
			w := httptest.NewRecorder()

			handler := NewVehicleHandler(mockSvc)
			handler.GetVehicleById(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != testCase.wantStatus {
				t.Errorf("expected %d, got %d", testCase.wantStatus, res.StatusCode)
			}
		})
	}
}

func TestDeleteVehicleHandler(t *testing.T) {
	const url = "/vehicles"

	tests := []struct {
		name       string
		id         string
		authInfo   *models.AuthInfo
		body       io.Reader
		wantStatus int
	}{
		{"no token", "1", nil, nil, http.StatusUnauthorized},
		{"ok", "1", testutils.GetAdminAuth(), nil, http.StatusNoContent},
		{"with body", "1", testutils.GetAdminAuth(), strings.NewReader(`{"invalid":"data"}`), http.StatusBadRequest},
		{"bad id in path", "invalid", testutils.GetAdminAuth(), nil, http.StatusBadRequest},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			mockSvc := &mockVehicleService{
				deleteByIdFunc: func(ctx context.Context, authInfo models.AuthInfo, id int) error {
					return nil
				},
			}

			req := testutils.CreateNewRequest(http.MethodPatch, url, testCase.authInfo, testCase.body)
			req = testutils.AddChiIdToContext(req, testCase.id)
			w := httptest.NewRecorder()

			handler := NewVehicleHandler(mockSvc)
			handler.DeleteVehicle(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != testCase.wantStatus {
				t.Errorf("expected %d, got %d", testCase.wantStatus, res.StatusCode)
			}
		})
	}
}
