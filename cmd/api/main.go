package main

import (
	"context"
	"log"

	config "github.com/OxytocinGroup/theca-backend/internal/config"
	db "github.com/OxytocinGroup/theca-backend/internal/db"
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
	database := db.ConnectDatabase(config).GetDB()
	var dependency di.DepsProvider
	if config.Environment == "dev" {
		deps := di.DevDeps{
			Config:    config,
			Db:        database,
			LogLogger: logger,
		}
		dependency = di.NewDevDeps(deps)
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
