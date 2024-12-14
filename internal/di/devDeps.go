package di

import (
	"github.com/OxytocinGroup/theca-backend/internal/config"
	db "github.com/OxytocinGroup/theca-backend/internal/db"
	"github.com/OxytocinGroup/theca-backend/internal/repository"
	"github.com/OxytocinGroup/theca-backend/internal/usecase"
	"gorm.io/gorm"
)

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
