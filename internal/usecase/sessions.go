package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/OxytocinGroup/theca-backend/internal/repository"
	"github.com/OxytocinGroup/theca-backend/pkg/logger"
)

type SessionUseCase interface {
	CreateSession(sessionID string, userID uint, expiresAt time.Time) error
	ValidateSession(sessionID string) (uint, error)
	DeleteSession(sessionID string) error
}

type sessionUseCase struct {
	sessionRepo repository.SessionRepository
	log         logger.Logger
}

func NewSessionUseCase(repo repository.SessionRepository, log logger.Logger) SessionUseCase {
	return &sessionUseCase{
		sessionRepo: repo,
		log:         log,
	}
}

func (suc *sessionUseCase) CreateSession(sessionID string, userID uint, expiresAt time.Time) error {
	return suc.sessionRepo.CreateSession(sessionID, userID, expiresAt)
}

func (suc *sessionUseCase) ValidateSession(sessionID string) (uint, error) {
	session, err := suc.sessionRepo.GetSessionByID(sessionID)
	if err != nil {
		suc.log.Error(context.Background(), "failed to get session by id", map[string]interface{}{
			"sessionID": sessionID,
		})
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
