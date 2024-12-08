package interfaces

import "time"

type SessionUseCase interface {
	CreateSession(sessionID string, userID uint, expiresAt time.Time) error
	ValidateSession(sessionID string) (uint, error)
	DeleteSession(sessionID string) error
}
