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

func TestGetVehiclesHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		url        string
		authInfo   *models.AuthInfo
		body       io.Reader
		wantStatus int
	}{
		{"no token", "GET", "/vehicles", nil, nil, http.StatusUnauthorized},
		{"ok", "GET", "/vehicles", testutils.GetAdminAuth(), nil, http.StatusOK},
		{"with body", "GET", "/vehicles", testutils.GetAdminAuth(), strings.NewReader(`{"invalid":"data"}`), http.StatusBadRequest},
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

			req := testutils.CreateNewRequest(testCase.method, testCase.url, testCase.authInfo, testCase.body)
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
