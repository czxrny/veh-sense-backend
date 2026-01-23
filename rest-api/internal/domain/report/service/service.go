package service

import (
	"context"
	"fmt"

	r "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/report/repository"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

type ReportService struct {
	repo *r.ReportRepository
}

func NewReportService(repo *r.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) FindAllReports(ctx context.Context, filter models.ReportFilter) ([]models.Report, error) {
	return s.repo.FindAll(ctx, filter)
}

func (s *ReportService) DeleteById(ctx context.Context, authInfo models.AuthInfo, id int) error {
	Report, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	isPrivateOwner := authInfo.OrganizationID == nil && Report.UserID == authInfo.UserID
	isOrgAdmin := Report.OrganizationID != nil && authInfo.OrganizationID != nil && Report.OrganizationID == authInfo.OrganizationID && authInfo.Role == "admin"

	if !isPrivateOwner && !isOrgAdmin && authInfo.Role != "root" {
		return fmt.Errorf("Error: User is unauthorized to delete the vehicle.")
	}

	return s.repo.DeleteById(ctx, id)
}
