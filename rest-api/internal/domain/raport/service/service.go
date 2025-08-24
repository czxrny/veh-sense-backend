package service

import (
	"context"
	"fmt"

	r "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/raport/repository"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

type RaportService struct {
	repo *r.RaportRepository
}

func NewRaportService(repo *r.RaportRepository) *RaportService {
	return &RaportService{repo: repo}
}

func (s *RaportService) FindAllReports(ctx context.Context, filter models.RaportFilter) ([]models.Raport, error) {
	return s.repo.FindAll(ctx, filter)
}

func (s *RaportService) DeleteById(ctx context.Context, authInfo models.AuthInfo, id int) error {
	raport, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	isPrivateOwner := authInfo.OrganizationID == nil && raport.UserID == authInfo.UserID
	isOrgAdmin := raport.OrganizationID != nil && authInfo.OrganizationID != nil && raport.OrganizationID == authInfo.OrganizationID && authInfo.Role == "admin"

	if !isPrivateOwner && !isOrgAdmin && authInfo.Role != "root" {
		return fmt.Errorf("Error: User is unauthorized to delete the vehicle.")
	}

	return s.repo.DeleteById(ctx, id)
}
