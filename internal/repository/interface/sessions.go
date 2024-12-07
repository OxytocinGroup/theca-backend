package interfaces

import (
	"time"

	"github.com/OxytocinGroup/theca-backend/internal/domain"
)

type SessionRepository interface {
	CreateSession(sessionID string, userID uint, expiresAt time.Time) error
	GetSessionByID(sessionID string) (domain.Session, error)
	DeleteSessionByID(sessionID string) error
}
