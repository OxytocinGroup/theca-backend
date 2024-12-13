package repository

import (
	domain "github.com/OxytocinGroup/theca-backend/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByEmail(email string) (domain.User, error)
	GetByUsername(username string) (domain.User, error)
	Create(user *domain.User) error
	EmailExists(email string) (bool, error)
	UsernameExists(username string) (bool, error)
	Update(user *domain.User) error
}

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) UserRepository {
	return &userDatabase{DB}
}

func (c *userDatabase) GetByEmail(email string) (domain.User, error) {
	var user domain.User
	err := c.DB.Model(&domain.User{}).Where("email = ?", email).First(&user).Error

	return user, err
}

func (c *userDatabase) GetByUsername(username string) (domain.User, error) {
	var user domain.User
	err := c.DB.Model(&domain.User{}).Where("username = ?", username).First(&user).Error
	return user, err
}
func (c *userDatabase) Create(user *domain.User) error {
	return c.DB.Model(&domain.User{}).Create(user).Error
}

func (c *userDatabase) EmailExists(email string) (bool, error) {
	var count int64
	err := c.DB.Model(&domain.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (c *userDatabase) UsernameExists(username string) (bool, error) {
	var count int64
	err := c.DB.Model(&domain.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (c *userDatabase) Update(user *domain.User) error {
	return c.DB.Model(&domain.User{}).Where("id = ?", user.ID).Save(user).Error
}
