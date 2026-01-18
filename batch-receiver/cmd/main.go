package main

import (
	"log"

	"github.com/czxrny/veh-sense-backend/batch-receiver/internal/app"
	"github.com/czxrny/veh-sense-backend/batch-receiver/internal/server"
)

func main() {
	App, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	defer app.GetSQLClient().Close()

	if err := server.InitializeAndStart(App); err != nil {
		log.Fatal("Couldn't start the Batch Receiver server: ", err)
	}
}
