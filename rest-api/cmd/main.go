package main

import (
	"log"

	server "github.com/czxrny/veh-sense-backend/rest-api/internal/server"
	"github.com/czxrny/veh-sense-backend/shared/database"
)

func main() {
	if err := database.ConnectToDatabase(); err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}
	defer database.GetSQLClient().Close()

	if err := server.InitializeAndStart(); err != nil {
		log.Fatal("Couldn't start the HTTP server: ", err)
	}
}
