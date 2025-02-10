package cron

import (
	"context"
	"fmt"
	"time"

	"github.com/OxytocinGroup/theca-backend/internal/config"
	"github.com/OxytocinGroup/theca-backend/internal/repository"
	"github.com/OxytocinGroup/theca-backend/pkg/logger"
	"github.com/go-co-op/gocron"
)

var (
	conf  *config.Config
	repos repository.SessionRepository
	logs  logger.Logger
)

func clearDB() {
	fmt.Println("in clear")
	logs.Info(context.Background(), "start cleaning session database", map[string]any{})

	sessions, err := repos.GetAllSessions()
	if err != nil {
		logs.Error(context.Background(), "cron (clear session db): error while getting sessions", map[string]any{"error": err})
	}

	for _, session := range sessions {
		if session.ExpiresAt.Before(time.Now()) {
			if err := repos.DeleteSessionByID(session.ID); err != nil {
				logs.Error(context.Background(), "cron (clear session db): error while deleting session", map[string]any{"error": err})
			}
		}
	}

}

func InitScheduler(cfg *config.Config, log logger.Logger, repo repository.SessionRepository) {
	conf = cfg
	repos = repo
	logs = log
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Error(context.Background(), "cron: failed to load timezone", map[string]any{"error": err})
	}

	scheduler := gocron.NewScheduler(location)

	scheduler.Every(1).Day().At(cfg.ClearTime).Do(clearDB)

	scheduler.StartAsync()
}
