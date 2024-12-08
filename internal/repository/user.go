package repository

import (
	"fmt"

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
	fmt.Println(user)
	return c.DB.Model(&domain.User{}).Where("id = ?", user.ID).Save(user).Error
}
