package usecase

import (
	"context"

	domain "github.com/OxytocinGroup/theca-backend/internal/domain"
	interfaces "github.com/OxytocinGroup/theca-backend/internal/repository/interface"
	services "github.com/OxytocinGroup/theca-backend/internal/usecase/interface"
)

type userUseCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUseCase(repo interfaces.UserRepository) services.UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}

func (c *userUseCase) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := c.userRepo.GetByEmail(ctx, email)
	return user, err
}

func (c *userUseCase) Create(ctx context.Context, user *domain.User) error {
	err := c.userRepo.Create(ctx, user)
	return err
}

func (c *userUseCase) EmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := c.userRepo.EmailExists(ctx, email)

	return exists, err
}

func (c *userUseCase) UsernameExists(ctx context.Context, username string) (bool, error) {
	exists, err := c.userRepo.UsernameExists(ctx, username)

	return exists, err
}
