package di

import (
	"github.com/OxytocinGroup/theca-backend/internal/config"
	"github.com/OxytocinGroup/theca-backend/internal/repository"
	"github.com/OxytocinGroup/theca-backend/internal/usecase"
	"github.com/OxytocinGroup/theca-backend/pkg/logger"
	"gorm.io/gorm"
)

type DevDeps struct {
	Config    config.Config
	Db        *gorm.DB
	LogLogger logger.Logger
}

func NewDevDeps(deps DevDeps) DepsProvider {
	return &DevDeps{
		Config:    deps.Config,
		Db:        deps.Db,
		LogLogger: deps.LogLogger,
	}
}

func (d *DevDeps) Database() *gorm.DB {
	return d.Db
}

func (d *DevDeps) SessionRepository() repository.SessionRepository {
	return repository.NewSessionRepository(d.Db)
}

func (d *DevDeps) UserRepository() repository.UserRepository {
	return repository.NewUserRepository(d.Db)
}

func (d *DevDeps) SessionUseCase(repo repository.SessionRepository, log logger.Logger) usecase.SessionUseCase {
	return usecase.NewSessionUseCase(repo, log)
}

func (d *DevDeps) UserUseCase(repo repository.UserRepository, log logger.Logger) usecase.UserUseCase {
	return usecase.NewUserUseCase(repo, log)
}

func (d *DevDeps) Logger() logger.Logger {
	return d.LogLogger
}
