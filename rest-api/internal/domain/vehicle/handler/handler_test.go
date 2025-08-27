package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/czxrny/veh-sense-backend/rest-api/internal/testutils"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

type mockVehicleService struct {
	findVehiclesFunc func(ctx context.Context, filter models.VehicleFilter) ([]models.Vehicle, error)
	addVehicleFunc   func(ctx context.Context, vehicle *models.Vehicle, authInfo models.AuthInfo) (*models.Vehicle, error)
	getByIdFunc      func(ctx context.Context, authInfo models.AuthInfo, id int) (*models.Vehicle, error)
	updateByIdFunc   func(ctx context.Context, authInfo models.AuthInfo, updatedVehicle *models.VehicleUpdate, id int) (*models.Vehicle, error)
	deleteByIdFunc   func(ctx context.Context, authInfo models.AuthInfo, id int) error
}

func (m *mockVehicleService) FindVehicles(ctx context.Context, filter models.VehicleFilter) ([]models.Vehicle, error) {
	return m.findVehiclesFunc(ctx, filter)
}

func (m *mockVehicleService) AddVehicle(ctx context.Context, vehicle *models.Vehicle, authInfo models.AuthInfo) (*models.Vehicle, error) {
	return m.addVehicleFunc(ctx, vehicle, authInfo)
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

func TestGetAllVehiclesNoToken(t *testing.T) {
	mockSvc := &mockVehicleService{
		findVehiclesFunc: func(ctx context.Context, filter models.VehicleFilter) ([]models.Vehicle, error) {
			return []models.Vehicle{}, nil
		},
	}

	req := testutils.CreateNewRequest(http.MethodGet, "/vehicles", nil, nil)
	w := httptest.NewRecorder()

	handler := NewVehicleHandler(mockSvc)
	handler.GetVehicles(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 OK, got %d", res.StatusCode)
	}
}

func TestGetAllVehiclesOK(t *testing.T) {
	mockSvc := &mockVehicleService{
		findVehiclesFunc: func(ctx context.Context, filter models.VehicleFilter) ([]models.Vehicle, error) {
			return []models.Vehicle{
				{ID: 1, Brand: "BMW", Model: "X5"},
				{ID: 2, Brand: "Audi", Model: "A4"},
			}, nil
		},
	}

	req := testutils.CreateNewRequest(http.MethodGet, "/vehicles", testutils.GetAdminAuth(), nil)
	w := httptest.NewRecorder()

	handler := NewVehicleHandler(mockSvc)
	handler.GetVehicles(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", res.StatusCode)
	}
}

func TestGetAllVehiclesWithBody(t *testing.T) {
	mockSvc := &mockVehicleService{
		findVehiclesFunc: func(ctx context.Context, filter models.VehicleFilter) ([]models.Vehicle, error) {
			return []models.Vehicle{
				{ID: 1, Brand: "BMW", Model: "X5"},
				{ID: 2, Brand: "Audi", Model: "A4"},
			}, nil
		},
	}

	req := testutils.CreateNewRequest(http.MethodGet, "/vehicles", testutils.GetAdminAuth(), strings.NewReader(`{"invalid":"data"}`))
	w := httptest.NewRecorder()

	handler := NewVehicleHandler(mockSvc)
	handler.GetVehicles(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 OK, got %d", res.StatusCode)
	}
}
