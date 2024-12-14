package di

import (
	"github.com/OxytocinGroup/theca-backend/internal/config"
	db "github.com/OxytocinGroup/theca-backend/internal/db"
	"github.com/OxytocinGroup/theca-backend/internal/repository"
	"github.com/OxytocinGroup/theca-backend/internal/usecase"
	"github.com/OxytocinGroup/theca-backend/pkg/logger"
	"gorm.io/gorm"
)

type DevDeps struct {
	config config.Config
	db     *gorm.DB
	logger logger.Logger
}

func NewDevDeps(cfg config.Config) DepsProvider {
	database := db.ConnectDatabase(cfg).GetDB()
	return &DevDeps{
		config: cfg,
		db:     database,
		logger: logger.NewLogrusLogger(cfg.LogLevel),
	}
}

func (d *DevDeps) Database() *gorm.DB {
	return d.db
}

func (d *DevDeps) SessionRepository() repository.SessionRepository {
	return repository.NewSessionRepository(d.db)
}

func (d *DevDeps) UserRepository() repository.UserRepository {
	return repository.NewUserRepository(d.db)
}

func (d *DevDeps) SessionUseCase(repo repository.SessionRepository, log logger.Logger) usecase.SessionUseCase {
	return usecase.NewSessionUseCase(repo, log)
}

func (d *DevDeps) UserUseCase(repo repository.UserRepository, log logger.Logger) usecase.UserUseCase {
	return usecase.NewUserUseCase(repo, log)
}

func (d *DevDeps) Logger() logger.Logger {
	return d.logger
}
