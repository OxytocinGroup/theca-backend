package main

import (
	"log"

	config "github.com/OxytocinGroup/theca-backend/internal/config"
	di "github.com/OxytocinGroup/theca-backend/internal/di"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	var dependency di.DepsProvider
	if config.Environment == "dev" {
		dependency = di.NewDevDeps(config)
	} else if config.Environment == "test" {
		// #TODO
	} else {
		// # TODO
	}

	server, diErr := di.InitializeAPI(dependency)
	if diErr != nil {
		log.Fatal("cannot start server: ", diErr)
	} else {
		server.Start()
	}
}
