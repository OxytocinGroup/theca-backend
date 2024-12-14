package di

import (
	http "github.com/OxytocinGroup/theca-backend/internal/api"
	"github.com/OxytocinGroup/theca-backend/internal/api/handler"

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

