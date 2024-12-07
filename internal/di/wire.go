//go:build wireinject
// +build wireinject

package di

import (
	http "github.com/OxytocinGroup/theca-backend/internal/api"
	handler "github.com/OxytocinGroup/theca-backend/internal/api/handler"
	config "github.com/OxytocinGroup/theca-backend/internal/config"
	db "github.com/OxytocinGroup/theca-backend/internal/db"
	repository "github.com/OxytocinGroup/theca-backend/internal/repository"
	usecase "github.com/OxytocinGroup/theca-backend/internal/usecase"
	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(db.ConnectDatabase, repository.NewUserRepository, usecase.NewUserUseCase, handler.NewUserHandler, http.NewServerHTTP)
	return &http.ServerHTTP{}, nil
}
