package interfaces

import (
	"github.com/OxytocinGroup/theca-backend/internal/domain"
)

type UserRepository interface {
	GetByEmail(email string) (domain.User, error)
	GetByUsername(username string) (domain.User, error)
	Create(user *domain.User) error
	EmailExists(email string) (bool, error)
	UsernameExists(username string) (bool, error)
	Update(user *domain.User) error
}
