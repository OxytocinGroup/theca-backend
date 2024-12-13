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

func (s *sessionUseCase) CreateSession(sessionID string, userID uint, expiresAt time.Time) error {
	return s.sessionRepo.CreateSession(sessionID, userID, expiresAt)
}

func (s *sessionUseCase) ValidateSession(sessionID string) (uint, error) {
	session, err := s.sessionRepo.GetSessionByID(sessionID)
	if err != nil {
		return 0, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return 0, errors.New("session expired")
	}
	return session.UserID, nil
}

func (s *sessionUseCase) DeleteSession(sessionID string) error {
	return s.sessionRepo.DeleteSessionByID(sessionID)
}
