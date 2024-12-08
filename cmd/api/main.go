package main

import (
	"log"

	config "github.com/OxytocinGroup/theca-backend/internal/config"
	di "github.com/OxytocinGroup/theca-backend/internal/di"
	"github.com/OxytocinGroup/theca-backend/pkg/logger"
)

func main() {
	err := logger.InitializeLogger("logs/app.log")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	server, diErr := di.InitializeAPI(config)
	if diErr != nil {
		log.Fatal("cannot start server: ", diErr)
	} else {
		server.Start()
	}
}
