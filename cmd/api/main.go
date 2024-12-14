package main

import (
	"context"
	"log"

	config "github.com/OxytocinGroup/theca-backend/internal/config"
	di "github.com/OxytocinGroup/theca-backend/internal/di"
	"github.com/OxytocinGroup/theca-backend/pkg/logger"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	logger := logger.NewLogrusLogger(config.LogLevel)
	logger.Info(context.Background(), "App started", nil)

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
