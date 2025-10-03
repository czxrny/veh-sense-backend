package main

import (
	"log"

	"github.com/czxrny/veh-sense-backend/batch-receiver/internal/app"
	"github.com/czxrny/veh-sense-backend/batch-receiver/internal/server"
)

func main() {
	if err := app.NewApp(); err != nil {
		log.Fatal(err)
	}
	if err := server.InitializeAndStart(); err != nil {
		log.Fatal("Couldn't start the Batch Receiver server: ", err)
	}
}
