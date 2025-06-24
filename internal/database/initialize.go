package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var databaseClient *sql.DB

func ConnectToDatabase() error {
	var err error
	databaseClient, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}

	if err = databaseClient.Ping(); err != nil {
		return err
	}
	fmt.Println("Connected to database!")
	return nil
}

func GetDatabaseClient() *sql.DB {
	return databaseClient
}
