package main

import (
	"log"
	"veh-sense-backend/internal/database"
	server "veh-sense-backend/internal/server/rest-api"
)

func main() {
	if err := database.ConnectToDatabase(); err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}
	defer database.GetDatabaseClient().Close()

	if err := server.InitializeAndStart(); err != nil {
		log.Fatal("Couldn't start the HTTP server: ", err)
	}
}
