package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/czxrny/veh-sense-backend/shared/models"

	vRepo "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/vehicle/repository"
	vServ "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/vehicle/service"

	oRepo "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/organization/repository"
	oServ "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/organization/service"

	rRepo "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/report/repository"
	rServ "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/report/service"

	uRepo "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/user/repository"
	uServ "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/user/service"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var databaseClient *gorm.DB

type App struct {
	VehicleService      vServ.VehicleService
	OrganizationService *oServ.OrganizationService
	ReportService       *rServ.ReportService
	UserService         *uServ.UserService
}

type repoList struct {
	Vehicle      *vRepo.VehicleRepository
	Organization *oRepo.OrganizationRepository
	Report       *rRepo.ReportRepository
	ReportData   *rRepo.ReportDataRepository
	UserAuth     *uRepo.UserAuthRepository
	UserInfo     *uRepo.UserInfoRepository
	RefreshKey   *uRepo.RefreshKeyRepository
}

func NewApp() (*App, error) {
	var err error
	databaseClient, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err = autoMigrate(databaseClient); err != nil {
		return nil, fmt.Errorf("migration failed: " + err.Error())
	}

	repoList := createRepos(databaseClient)

	return &App{
		VehicleService:      vServ.NewVehicleService(repoList.Vehicle),
		OrganizationService: oServ.NewOrganizationService(repoList.Organization),
		ReportService:       rServ.NewReportService(repoList.Report, repoList.ReportData),
		UserService:         uServ.NewUserService(repoList.UserAuth, repoList.UserInfo, repoList.RefreshKey),
	}, nil
}

func autoMigrate(databaseClient *gorm.DB) error {
	if err := databaseClient.AutoMigrate(&models.Vehicle{}); err != nil {
		return err
	}

	if err := databaseClient.AutoMigrate(&models.Organization{}); err != nil {
		return err
	}

	if err := databaseClient.AutoMigrate(&models.UserAuth{}); err != nil {
		return err
	}

	if err := databaseClient.AutoMigrate(&models.UserInfo{}); err != nil {
		return err
	}

	if err := databaseClient.AutoMigrate(&models.RefreshInfo{}); err != nil {
		return err
	}

	return databaseClient.AutoMigrate(&models.Report{})
}

func GetDatabaseClient() *gorm.DB {
	return databaseClient
}

func GetSQLClient() *sql.DB {
	sqlDB, err := databaseClient.DB()
	if err != nil {
		log.Fatal(err)
	}
	return sqlDB
}

func createRepos(databaseClient *gorm.DB) repoList {
	return repoList{
		Vehicle:      vRepo.NewVehicleRepository(databaseClient),
		Organization: oRepo.NewOrganizationRepository(databaseClient),
		Report:       rRepo.NewReportRepository(databaseClient),
		ReportData:   rRepo.NewReportDataRepository(databaseClient),
		UserAuth:     uRepo.NewUserAuthRepository(databaseClient),
		UserInfo:     uRepo.NewUserInfoRepository(databaseClient),
		RefreshKey:   uRepo.NewRefreshKeyRepository(databaseClient),
	}

}
