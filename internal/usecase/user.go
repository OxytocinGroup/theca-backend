package usecase

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	utils "github.com/OxytocinGroup/theca-backend/internal/api/utils/email"
	"github.com/OxytocinGroup/theca-backend/internal/config"
	"github.com/OxytocinGroup/theca-backend/internal/domain"
	interfaces "github.com/OxytocinGroup/theca-backend/internal/repository/interface"
	services "github.com/OxytocinGroup/theca-backend/internal/usecase/interface"
	"github.com/OxytocinGroup/theca-backend/pkg"
)

type userUseCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUseCase(repo interfaces.UserRepository) services.UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}

func (c *userUseCase) Register(email, password, username string) pkg.Response {
	var user domain.User
	user.Email = email
	user.Password = password
	user.Username = username

	emailExists, err := c.userRepo.EmailExists(user.Email)
	if err != nil {
		return pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check email existence",
		}
	}
	if emailExists {
		return pkg.Response{
			Code:    http.StatusConflict,
			Message: "Email already exists",
		}
	}

	usernameExists, err := c.userRepo.UsernameExists(user.Username)
	if err != nil {
		return pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check username existence",
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

	mail := utils.Mail{Email: user.Email, Code: user.VerificationCode, Username: user.Username}
	config, err := config.LoadConfig()
	if err != nil {
		return pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to load config",
		}
	}
	if err := mail.SendVerificationEmail(config, user.Email, user.VerificationCode, user.Username); err != nil {
		fmt.Println(err)
		// #TODO LOGGER
		return pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to send verification email",
		}
	}
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
