package repository

import (
	"time"

	domain "github.com/OxytocinGroup/theca-backend/internal/domain"
	"gorm.io/gorm"
)

type SessionRepository interface {
	CreateSession(sessionID string, userID uint, expiresAt time.Time) error
	GetSessionByID(sessionID string) (domain.Session, error)
	DeleteSessionByID(sessionID string) error
	DeleteAllSessions(userID uint) error
}

type sessionDatabase struct {
	DB *gorm.DB
}

func NewSessionRepository(DB *gorm.DB) SessionRepository {
	return &sessionDatabase{DB}
}

func (sdb *sessionDatabase) CreateSession(sessionID string, userID uint, expiresAt time.Time) error {
	session := &domain.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
	return sdb.DB.Model(&domain.Session{}).Create(&session).Error
}

func (sdb *sessionDatabase) GetSessionByID(sessionID string) (domain.Session, error) {
	var session domain.Session
	err := sdb.DB.Model(&domain.Session{}).Where("id = ?", sessionID).First(&session).Error
	return session, err
}

func (sdb *sessionDatabase) DeleteSessionByID(sessionID string) error {
	return sdb.DB.Model(&domain.Session{}).Where("id = ?", sessionID).Delete(&domain.Session{}).Error
}

func (sdb *sessionDatabase) DeleteAllSessions(userID uint) error {
	return sdb.DB.Model(&domain.Session{}).Where("user_id = ?", userID).Delete(&domain.Session{}).Error
}
