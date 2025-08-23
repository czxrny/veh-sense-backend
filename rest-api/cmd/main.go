package main

import (
	"log"

	database "github.com/czxrny/veh-sense-backend/rest-api/internal/app"
	server "github.com/czxrny/veh-sense-backend/rest-api/internal/server"
)

func main() {
	App, err := database.NewApp()
	if err != nil {
		log.Fatal("Cannot initialize the database: ", err)
	}
	defer database.GetSQLClient().Close()

	if err := server.InitializeAndStart(App); err != nil {
		log.Fatal("Couldn't start the HTTP server: ", err)
	}
}
