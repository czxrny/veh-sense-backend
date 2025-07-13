package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/czxrny/veh-sense-backend/shared/models"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var databaseClient *gorm.DB

func ConnectToDatabase() error {
	var err error
	databaseClient, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := databaseClient.AutoMigrate(&models.Vehicle{}); err != nil {
		return fmt.Errorf("migration failed: " + err.Error())
	}

	if err := databaseClient.AutoMigrate(&models.Organization{}); err != nil {
		return fmt.Errorf("migration failed: " + err.Error())
	}

	if err := databaseClient.AutoMigrate(&models.UserAuth{}); err != nil {
		return fmt.Errorf("migration failed: " + err.Error())
	}

	if err := databaseClient.AutoMigrate(&models.UserInfo{}); err != nil {
		return fmt.Errorf("migration failed: " + err.Error())
	}

	if err := databaseClient.AutoMigrate(&models.Raport{}); err != nil {
		return fmt.Errorf("migration failed: " + err.Error())
	}

	return nil
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
