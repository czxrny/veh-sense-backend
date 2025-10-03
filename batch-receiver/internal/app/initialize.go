package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/czxrny/veh-sense-backend/shared/models"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var databaseClient *gorm.DB

func NewApp() error {
	var err error
	databaseClient, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
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
