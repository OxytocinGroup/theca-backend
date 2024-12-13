package usecase

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"

	"golang.org/x/crypto/bcrypt"

	utils "github.com/OxytocinGroup/theca-backend/internal/api/utils/email"
	"github.com/OxytocinGroup/theca-backend/internal/config"
	"github.com/OxytocinGroup/theca-backend/internal/domain"
	"github.com/OxytocinGroup/theca-backend/internal/repository"
	"github.com/OxytocinGroup/theca-backend/pkg"
)

type UserUseCase interface {
	Register(email, password, username string) pkg.Response
	VerifyEmail(email, code string) pkg.Response
	Auth(username, password string) (*domain.User, error)
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}

func (c *userUseCase) Register(email, password, username string) pkg.Response {
	var user domain.User
	user.Email = email
	user.Password = password
	user.Username = username

	var emailExists, usernameExists bool
	var emailError, usernameError error

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		emailExists, emailError = c.userRepo.EmailExists(user.Email)
	}()

	go func() {
		defer wg.Done()
		usernameExists, usernameError = c.userRepo.UsernameExists(user.Username)
	}()

	if emailError != nil || usernameError != nil {
		return pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check email or username existence",
		}
	}
	if emailExists {
		return pkg.Response{
			Code:    http.StatusConflict,
			Message: "Email already exists",
		}
	}
	if usernameExists {
		return pkg.Response{
			Code:    http.StatusConflict,
			Message: "Username already exists",
		}
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to hash password",
		}
	}

	user.Password = string(hashPass)
	user.VerificationCode = strconv.Itoa(rand.Intn(900000) + 100000)
	user.IsVerified = false

	if err := c.userRepo.Create(&user); err != nil {
		return pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
		}
	}

	configCh := make(chan *config.Config)
	mailErrorCh := make(chan error)

	go func() {
		config, err := config.LoadConfig()
		if err != nil {
			configCh <- nil
			return
		}
		configCh <- &config
	}()

	go func() {
		config := <-configCh
		if config == nil {
			mailErrorCh <- fmt.Errorf("failed to load config")
			return
		}
		mail := utils.Mail{Email: user.Email, Code: user.VerificationCode, Username: user.Username}
		err := mail.SendVerificationEmail(config, user.Email, user.VerificationCode, user.Username)
		mailErrorCh <- err
	}()

	return pkg.Response{
		Code:    http.StatusCreated,
		Message: "User registered successfully",
	}
}

func (c *userUseCase) VerifyEmail(email, code string) pkg.Response {
	user, err := c.userRepo.GetByEmail(email)
	if err != nil {
		return pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get user by email",
		}
	}

	if user.VerificationCode != code {
		return pkg.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid verification code",
		}
	}

	user.IsVerified = true
	user.VerificationCode = ""

	if err := c.userRepo.Update(&user); err != nil {
		return pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update user",
		}
	}

	return pkg.Response{
		Code:    http.StatusOK,
		Message: "Email verified successfully",
	}
}

func (c *userUseCase) Auth(username, password string) (*domain.User, error) {
	user, err := c.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return &user, err
}
