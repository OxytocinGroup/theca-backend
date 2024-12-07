package repository

import (
	"context"

	domain "github.com/OxytocinGroup/theca-backend/internal/domain"
	interfaces "github.com/OxytocinGroup/theca-backend/internal/repository/interface"
	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabase{DB}
}

func (c *userDatabase) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	err := c.DB.Where("email = ?", email).First(&user).Error

	return user, err
}

func (c *userDatabase) Create(ctx context.Context, user *domain.User) error {
	return c.DB.Create(user).Error
}

func (c *userDatabase) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := c.DB.Model(&domain.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (c *userDatabase) UsernameExists(ctx context.Context, username string) (bool, error) {
	var count int64
	err := c.DB.Model(&domain.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
