package database

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

	rRepo "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/raport/repository"
	rServ "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/raport/service"

	uRepo "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/user/repository"
	uServ "github.com/czxrny/veh-sense-backend/rest-api/internal/domain/user/service"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var databaseClient *gorm.DB

type App struct {
	VehicleService      *vServ.VehicleService
	OrganizationService *oServ.OrganizationService
	RaportService       *rServ.RaportService
	UserService         *uServ.UserService
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

	VehicleRepo := vRepo.NewVehicleRepository(databaseClient)
	OrganizationRepo := oRepo.NewOrganizationRepository(databaseClient)
	RaportRepo := rRepo.NewRaportRepository(databaseClient)
	UserAuthRepository := uRepo.NewUserAuthRepository(databaseClient)
	UserInfoRepository := uRepo.NewUserInfoRepository(databaseClient)

	VehicleService := vServ.NewVehicleService(VehicleRepo)
	OrganizationService := oServ.NewOrganizationService(OrganizationRepo)
	RaportService := rServ.NewRaportService(RaportRepo)
	UserService := uServ.NewUserService(UserAuthRepository, UserInfoRepository)

	return &App{
		VehicleService:      VehicleService,
		OrganizationService: OrganizationService,
		RaportService:       RaportService,
		UserService:         UserService,
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

	return databaseClient.AutoMigrate(&models.Raport{})
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
