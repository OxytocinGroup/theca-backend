package repository

import (
	"time"

	domain "github.com/OxytocinGroup/theca-backend/internal/domain"
	interfaces "github.com/OxytocinGroup/theca-backend/internal/repository/interface"
	"gorm.io/gorm"
)

type sessionDatabase struct {
	DB *gorm.DB
}

func NewSessionRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabase{DB}
}

func (s *sessionDatabase) CreateSession(sessionID string, userID uint, expiresAt time.Time) error {
	session := &domain.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
	return s.DB.Create(&session).Error
}

func (s *sessionDatabase) GetSessionByID(sessionID string) (domain.Session, error) {
	var session domain.Session
	err := s.DB.Model(&domain.Session{}).Where("ID = ?", sessionID).First(&session).Error
	return session, err
}

func (s *sessionDatabase) DeleteSessionByID(sessionID string) error {
	return s.DB.Model(&domain.Session{}).Where("ID = ?", sessionID).Delete(&domain.Session{}).Error
}
