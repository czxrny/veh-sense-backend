package service

import (
	"context"
	"fmt"

	r "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/vehicle/repository"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

type VehicleService struct {
	repo *r.VehicleRepository
}

func NewVehicleService(repo *r.VehicleRepository) *VehicleService {
	return &VehicleService{repo: repo}
}

func (s *VehicleService) FindVehicles(ctx context.Context, filter models.VehicleFilter) ([]models.Vehicle, error) {
	return s.repo.FindAll(ctx, filter)
}

func (s *VehicleService) AddVehicle(ctx context.Context, vehicle *models.Vehicle, authInfo models.AuthInfo) (*models.Vehicle, error) {
	switch authInfo.Role {
	case "user":
		if authInfo.OrganizationID != nil {
			return nil, fmt.Errorf("Error: Corporate user is unauthorized to add new vehicles to fleet. Login as an admin to proceed.")
		}
		vehicle.OwnerID = &authInfo.UserID
		vehicle.OrganizationID = authInfo.OrganizationID

	// Admin sets the vehicle organization_id automatically, and can pass the owner_id to specify the user that will be the owner
	case "admin":
		vehicle.OrganizationID = authInfo.OrganizationID

	case "root":
		if vehicle.OwnerID == nil || vehicle.OrganizationID == nil {
			return nil, fmt.Errorf("Error: Bad Request: Please specify either the organization_id or owner_id to proceed")
		}
	}

	vehicle.ID = 0

	err := s.repo.Add(ctx, vehicle)

	return vehicle, err
}

func (s *VehicleService) GetById(ctx context.Context, authInfo models.AuthInfo, id int) (*models.Vehicle, error) {
	vehicle, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	isOwner := vehicle.OwnerID != nil && *vehicle.OwnerID == authInfo.UserID
	isShared := vehicle.OrganizationID != nil && authInfo.OrganizationID != nil && *vehicle.OrganizationID == *authInfo.OrganizationID && vehicle.OwnerID == nil
	isOrgAdmin := vehicle.OrganizationID != nil && authInfo.OrganizationID != nil && *vehicle.OrganizationID == *authInfo.OrganizationID && authInfo.Role == "admin"

	if !isOwner && !isShared && !isOrgAdmin && authInfo.Role != "root" {
		return nil, fmt.Errorf("Error: User is unauthorized to view the vehicle.")
	}

	return vehicle, nil
}

func (s *VehicleService) UpdateById(ctx context.Context, authInfo models.AuthInfo, updatedVehicle *models.VehicleUpdate, id int) (*models.Vehicle, error) {
	vehicle, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	isPrivateOwner := authInfo.OrganizationID == nil && vehicle.OwnerID != nil && *vehicle.OwnerID == authInfo.UserID
	isOrgAdmin := vehicle.OrganizationID != nil && authInfo.OrganizationID != nil && *vehicle.OrganizationID == *authInfo.OrganizationID && authInfo.Role == "admin"

	if !isPrivateOwner && !isOrgAdmin && authInfo.Role != "root" {
		return nil, fmt.Errorf("Error: User is unauthorized to edit the vehicle.")
	}

	if err := s.repo.UpdatePartial(ctx, id, updatedVehicle); err != nil {
		return nil, err
	}

	// returning updated object
	return s.repo.GetByID(ctx, id)
}

func (s *VehicleService) DeleteById(ctx context.Context, authInfo models.AuthInfo, id int) error {
	vehicle, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	isPrivateOwner := authInfo.OrganizationID == nil && vehicle.OwnerID != nil && *vehicle.OwnerID == authInfo.UserID
	isOrgAdmin := vehicle.OrganizationID != nil && authInfo.OrganizationID != nil && *vehicle.OrganizationID == *authInfo.OrganizationID && authInfo.Role == "admin"

	if !isPrivateOwner && !isOrgAdmin && authInfo.Role != "root" {
		return fmt.Errorf("Error: User is unauthorized to delete the vehicle.")
	}

	return s.repo.Delete(ctx, vehicle)
}
