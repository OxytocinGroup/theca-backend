package usecase

import (
	"errors"
	"time"

	"github.com/OxytocinGroup/theca-backend/internal/repository"
)

type SessionUseCase interface {
	CreateSession(sessionID string, userID uint, expiresAt time.Time) error
	ValidateSession(sessionID string) (uint, error)
	DeleteSession(sessionID string) error
}

type sessionUseCase struct {
	sessionRepo repository.SessionRepository
}

func NewSessionUseCase(repo repository.SessionRepository) SessionUseCase {
	return &sessionUseCase{
		sessionRepo: repo,
	}
}

func (suc *sessionUseCase) CreateSession(sessionID string, userID uint, expiresAt time.Time) error {
	return suc.sessionRepo.CreateSession(sessionID, userID, expiresAt)
}

func (suc *sessionUseCase) ValidateSession(sessionID string) (uint, error) {
	session, err := suc.sessionRepo.GetSessionByID(sessionID)
	if err != nil {
		return 0, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return 0, errors.New("session expired")
	}
	return session.UserID, nil
}

func (suc *sessionUseCase) DeleteSession(sessionID string) error {
	return suc.sessionRepo.DeleteSessionByID(sessionID)
}
