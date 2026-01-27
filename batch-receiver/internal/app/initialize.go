package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/czxrny/veh-sense-backend/batch-receiver/internal/domain/upload/repository"
	s "github.com/czxrny/veh-sense-backend/batch-receiver/internal/domain/upload/service"
	"github.com/czxrny/veh-sense-backend/shared/models"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var databaseClient *gorm.DB

type App struct {
	Service s.Service
}

func NewApp() (*App, error) {
	var err error
	databaseClient, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// just to be sure report exists - if rest api did not run auto migrate
	err = databaseClient.AutoMigrate(&models.Report{})
	if err != nil {
		return nil, err
	}

	err = databaseClient.AutoMigrate(&models.RawRideRecord{})
	if err != nil {
		return nil, err
	}

	reportDataRepo := repository.NewReportDataRepository(databaseClient)
	reportRepo := repository.NewReportRepository(databaseClient)
	userRepo := repository.NewUserInfoRepository(databaseClient)

	return &App{
		Service: *s.NewService(reportRepo, reportDataRepo, userRepo),
	}, nil
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
