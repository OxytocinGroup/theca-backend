package di

import (
	http "github.com/OxytocinGroup/theca-backend/internal/api"
	"github.com/OxytocinGroup/theca-backend/internal/api/handler"
	"github.com/OxytocinGroup/theca-backend/internal/config"
	db "github.com/OxytocinGroup/theca-backend/internal/db"
	"github.com/OxytocinGroup/theca-backend/internal/repository"
	"github.com/OxytocinGroup/theca-backend/internal/usecase"
	"gorm.io/gorm"
)

type DepsProvider interface {
	Database() *gorm.DB
	UserRepository() repository.UserRepository
	SessionRepository() repository.SessionRepository
	UserUseCase(repository.UserRepository) usecase.UserUseCase
	SessionUseCase(repository.SessionRepository) usecase.SessionUseCase
}

func InitializeAPI(provider DepsProvider) (*http.ServerHTTP, error) {
	userRepo := provider.UserRepository()
	sessionRepo := provider.SessionRepository()

	userUC := provider.UserUseCase(userRepo)
	sessionUC := provider.SessionUseCase(sessionRepo)

	userHandler := handler.NewUserHandler(userUC, sessionUC)
	return http.NewServerHTTP(userHandler), nil
}

type DevDeps struct {
	config config.Config
	db     *gorm.DB
}

func NewDevDeps(cfg config.Config) DepsProvider {
	database := db.ConnectDatabase(cfg).GetDB()
	return &DevDeps{config: cfg, db: database}

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

func (d *DevDeps) SessionUseCase(repo repository.SessionRepository) usecase.SessionUseCase {
	return usecase.NewSessionUseCase(repo)
}

func (d *DevDeps) UserUseCase(repo repository.UserRepository) usecase.UserUseCase {
	return usecase.NewUserUseCase(repo)
}
