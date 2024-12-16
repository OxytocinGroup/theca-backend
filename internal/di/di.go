package di

import (
	http "github.com/OxytocinGroup/theca-backend/internal/api"
	"github.com/OxytocinGroup/theca-backend/internal/api/handler"
	"github.com/OxytocinGroup/theca-backend/internal/repository"
	"github.com/OxytocinGroup/theca-backend/internal/usecase"
	"github.com/OxytocinGroup/theca-backend/pkg/logger"
	"gorm.io/gorm"
)

type DepsProvider interface {
	Database() *gorm.DB
	UserRepository() repository.UserRepository
	SessionRepository() repository.SessionRepository
	BookmarkRepository() repository.BookmarkRepository

	UserUseCase(repository.UserRepository, repository.SessionRepository, logger.Logger) usecase.UserUseCase
	SessionUseCase(repository.SessionRepository, logger.Logger) usecase.SessionUseCase
	BookmarkUseCase(repository.BookmarkRepository, logger.Logger) usecase.BookmarkUseCase

	Logger() logger.Logger
}

func InitializeAPI(provider DepsProvider) (*http.ServerHTTP, error) {
	log := provider.Logger()

	userRepo := provider.UserRepository()
	sessionRepo := provider.SessionRepository()
	bookmarkRepo := provider.BookmarkRepository()

	userUC := provider.UserUseCase(userRepo, sessionRepo, log)
	sessionUC := provider.SessionUseCase(sessionRepo, log)
	bookmarkUC := provider.BookmarkUseCase(bookmarkRepo, log)

	userHandler := handler.NewUserHandler(userUC, sessionUC, log)
	bookmarkHandler := handler.NewBookmarkHandler(bookmarkUC, log)
	return http.NewServerHTTP(userHandler, bookmarkHandler), nil
}
