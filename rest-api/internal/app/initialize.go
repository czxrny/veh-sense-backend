package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/czxrny/veh-sense-backend/shared/models"

	vr "github.com/czxrny/veh-sense-backend/rest-api/internal/repositories/vehicle"
	vs "github.com/czxrny/veh-sense-backend/rest-api/internal/services/vehicle"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var databaseClient *gorm.DB

type App struct {
	VehicleService *vs.VehicleService
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

	VehicleRepo := vr.NewVehicleRepository(databaseClient)

	VehicleService := vs.NewVehicleService(VehicleRepo)

	return &App{
		VehicleService: VehicleService,
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
