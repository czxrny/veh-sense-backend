package main

import (
	server "veh-sense-backend/internal/server/rest-api"
)

func main() {
	if err := server.InitializeAndStart(); err != nil {
		panic(err)
	}
}
