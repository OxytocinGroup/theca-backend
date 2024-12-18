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
	GetByID(id string) (domain.User, error)
	CheckVerificationStatus(userID uint) (bool, error)
	GetByToken(token string) (domain.User, error)
}

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) UserRepository {
	return &userDatabase{DB}
}

func (udb *userDatabase) GetByEmail(email string) (domain.User, error) {
	var user domain.User
	err := udb.DB.Model(&domain.User{}).Where("email = ?", email).First(&user).Error

	return user, err
}

func (udb *userDatabase) GetByUsername(username string) (domain.User, error) {
	var user domain.User
	err := udb.DB.Model(&domain.User{}).Where("username = ?", username).First(&user).Error
	return user, err
}
func (udb *userDatabase) Create(user *domain.User) error {
	return udb.DB.Model(&domain.User{}).Create(user).Error
}

func (udb *userDatabase) EmailExists(email string) (bool, error) {
	var count int64
	err := udb.DB.Model(&domain.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (udb *userDatabase) UsernameExists(username string) (bool, error) {
	var count int64
	err := udb.DB.Model(&domain.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (udb *userDatabase) Update(user *domain.User) error {
	return udb.DB.Model(&domain.User{}).Where("id = ?", user.ID).Save(user).Error
}

func (udb *userDatabase) GetByID(id string) (domain.User, error) {
	var user domain.User
	err := udb.DB.Model(&domain.User{}).Where("id = ?", id).First(&user).Error
	return user, err
}

func (udb *userDatabase) CheckVerificationStatus(userID uint) (bool, error) {
	var status bool
	err := udb.DB.Model(&domain.User{}).Where("id = ?", userID).Pluck("is_verified", &status).Error
	return status, err
}

func (udb *userDatabase) GetByToken(token string) (domain.User, error) {
	var user domain.User
	err := udb.DB.Model(&domain.User{}).Where("reset_token = ?", token).First(&user).Error

	return user, err
}
