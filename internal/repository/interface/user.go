package interfaces

import (
	"context"

	"github.com/OxytocinGroup/theca-backend/internal/domain"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	EmailExists(ctx context.Context, email string) (bool, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
}
