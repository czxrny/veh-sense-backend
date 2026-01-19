package app

import (
	"database/sql"
	"log"
	"os"

	r "github.com/czxrny/veh-sense-backend/batch-receiver/internal/domain/raport/repository"
	s "github.com/czxrny/veh-sense-backend/batch-receiver/internal/domain/raport/service"
	"github.com/czxrny/veh-sense-backend/batch-receiver/internal/model"
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

	// just to be sure raport exists - if rest api did not run auto migrate
	err = databaseClient.AutoMigrate(&models.Raport{})
	if err != nil {
		return nil, err
	}

	err = databaseClient.AutoMigrate(&model.RawRideRecord{})
	if err != nil {
		return nil, err
	}

	raportRepo := r.NewRaportRepository(databaseClient)
	raportDataRepo := r.NewRaportDataRepository(databaseClient)
	userRepo := r.NewUserInfoRepository(databaseClient)

	return &App{
		Service: *s.NewService(raportRepo, raportDataRepo, userRepo),
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
