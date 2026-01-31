package service

import (
	"context"
	"encoding/base64"
	"fmt"

	r "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/report/repository"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

type ReportService struct {
	reportRepo *r.ReportRepository
	recordRepo *r.ReportDataRepository
}

func NewReportService(repo *r.ReportRepository, recordRepo *r.ReportDataRepository) *ReportService {
	return &ReportService{reportRepo: repo, recordRepo: recordRepo}
}

func (s *ReportService) FindAllReports(ctx context.Context, filter models.ReportFilter) ([]models.Report, error) {
	return s.reportRepo.FindAll(ctx, filter)
}

func (s *ReportService) FindAllReportsOrganization(ctx context.Context, filter models.ReportFilter) ([]models.AdminReport, error) {
	return s.reportRepo.FindAllAdmin(ctx, filter)
}

func (s *ReportService) FindById(ctx context.Context, authInfo models.AuthInfo, id int) (*models.Report, error) {
	report, err := s.reportRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	isOwner := report.UserID == authInfo.UserID
	isOrgAdmin := report.OrganizationID != nil && authInfo.OrganizationID != nil && *report.OrganizationID == *authInfo.OrganizationID && authInfo.Role == "admin"

	if !isOwner && !isOrgAdmin {
		return nil, fmt.Errorf("Error: User is unauthorized to get the report.")
	}

	return report, nil
}

func (s *ReportService) DeleteById(ctx context.Context, authInfo models.AuthInfo, id int) error {
	Report, err := s.reportRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	isPrivateOwner := authInfo.OrganizationID == nil && Report.UserID == authInfo.UserID
	isOrgAdmin := Report.OrganizationID != nil && authInfo.OrganizationID != nil && Report.OrganizationID == authInfo.OrganizationID && authInfo.Role == "admin"

	if !isPrivateOwner && !isOrgAdmin && authInfo.Role != "root" {
		return fmt.Errorf("Error: User is unauthorized to delete the vehicle.")
	}

	return s.reportRepo.DeleteById(ctx, id)
}

func (s *ReportService) GetRideData(ctx context.Context, authInfo models.AuthInfo, id int) (*models.RideRecord, error) {
	report, err := s.reportRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Internal error while fetching data from db" + err.Error())
	}

	isOwner := report.UserID == authInfo.UserID
	isOrgAdmin := report.OrganizationID != nil && authInfo.OrganizationID != nil && *report.OrganizationID == *authInfo.OrganizationID && authInfo.Role == "admin"

	if !isOwner && !isOrgAdmin {
		return nil, fmt.Errorf("Error: User is unauthorized to access the data.")
	}

	raw, err := s.recordRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Internal error while fetching data from db" + err.Error())
	}

	return &models.RideRecord{
		ReportID:  raw.ReportID,
		Data:      base64.StdEncoding.EncodeToString(raw.Data),
		EventData: base64.StdEncoding.EncodeToString(raw.EventData),
	}, nil
}
