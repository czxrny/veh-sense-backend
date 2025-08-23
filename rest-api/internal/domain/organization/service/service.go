package service

import (
	"context"
	"fmt"

	r "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/organization/repository"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

type OrganizationService struct {
	repo *r.OrganizationRepository
}

func NewOrganizationService(repo *r.OrganizationRepository) *OrganizationService {
	return &OrganizationService{repo: repo}
}

func (s *OrganizationService) FindOrganizations(ctx context.Context, authInfo models.AuthInfo, filter models.OrganizationFilter) ([]models.Organization, error) {
	if authInfo.Role != "root" {
		return nil, fmt.Errorf("Error: User is unauthorized to add the asset.")
	}

	return s.repo.FindAll(ctx, filter)
}

func (s *OrganizationService) AddOrganization(ctx context.Context, organization *models.Organization, authInfo models.AuthInfo) (*models.Organization, error) {
	if authInfo.Role != "root" {
		return nil, fmt.Errorf("Error: User is unauthorized to add the asset.")
	}

	organization.ID = 0
	err := s.repo.Add(ctx, organization)

	return organization, err
}

func (s *OrganizationService) GetMyOrganizationInfo(ctx context.Context, authInfo models.AuthInfo) (*models.Organization, error) {
	if authInfo.OrganizationID == nil {
		return nil, fmt.Errorf("Error: User does not have an organization.")
	}

	return s.repo.GetByID(ctx, *authInfo.OrganizationID)
}

func (s *OrganizationService) UpdateMyOrganization(ctx context.Context, authInfo models.AuthInfo, organizationUpdate *models.OrganizationUpdate) (*models.Organization, error) {
	isOrgAdmin := authInfo.OrganizationID != nil && authInfo.Role == "admin"
	if !isOrgAdmin && authInfo.Role != "root" {
		return nil, fmt.Errorf("Error: User is unauthorized to edit the organization.")
	}

	_, err := s.repo.GetByID(ctx, *authInfo.OrganizationID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdatePartial(ctx, *authInfo.OrganizationID, organizationUpdate); err != nil {
		return nil, err
	}

	// returning updated object
	return s.repo.GetByID(ctx, *authInfo.OrganizationID)
}

func (s *OrganizationService) DeleteById(ctx context.Context, authInfo models.AuthInfo, id int) error {
	organization, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if authInfo.Role != "root" {
		return fmt.Errorf("Error: User is unauthorized to add the asset.")
	}

	return s.repo.Delete(ctx, organization)
}
