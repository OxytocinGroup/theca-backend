package usecase

import (
	"context"
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
	"github.com/OxytocinGroup/theca-backend/pkg/logger"
)

type UserUseCase interface {
	Register(email, password, username string) pkg.Response
	VerifyEmail(email, code string) pkg.Response
	Auth(username, password string) (*domain.User, pkg.Response)
	ChangePass(userID string, newPassword string) pkg.Response
}

type userUseCase struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	log         logger.Logger
}

func NewUserUseCase(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, log logger.Logger) UserUseCase {
	return &userUseCase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		log:         log,
	}
}

func (uuc *userUseCase) Register(email, password, username string) pkg.Response {
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
		emailExists, emailError = uuc.userRepo.EmailExists(user.Email)
	}()

	go func() {
		defer wg.Done()
		usernameExists, usernameError = uuc.userRepo.UsernameExists(user.Username)
	}()

	wg.Wait()

	if emailError != nil || usernameError != nil {
		uuc.log.Error(context.Background(), "failed to check email or username existence", map[string]interface{}{
			"email error":    emailError,
			"username error": usernameError,
		})
		return pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check email or username existence",
		}
	}
	if emailExists {
		uuc.log.Info(context.Background(), "email already exists", map[string]interface{}{
			"code": 409,
		})
		return pkg.Response{
			Code:    http.StatusConflict,
			Message: "Email already exists",
		}
	}
	if usernameExists {
		uuc.log.Info(context.Background(), "username already exists", map[string]interface{}{
			"code": 409,
		})
		return pkg.Response{
			Code:    http.StatusConflict,
			Message: "Username already exists",
		}
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		uuc.log.Error(context.Background(), "failed to hash password", map[string]interface{}{
			"error": err,
		})
		return pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to hash password",
		}
	}

	user.Password = string(hashPass)
	user.VerificationCode = strconv.Itoa(rand.Intn(900000) + 100000)
	user.IsVerified = false

	if err := uuc.userRepo.Create(&user); err != nil {
		uuc.log.Error(context.Background(), "failed to create user", map[string]interface{}{
			"user":  user,
			"error": err,
		})
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
			uuc.log.Error(context.Background(), "failed to load config", map[string]interface{}{
				"error": err,
			})
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
		uuc.log.Error(context.Background(), "failed to send verification email", map[string]interface{}{
			"user_id": user.ID,
			"error":   err,
		})
	}()

	return pkg.Response{
		Code:    http.StatusCreated,
		Message: "User registered successfully",
	}
}

func (uuc *userUseCase) VerifyEmail(email, code string) pkg.Response {
	user, err := uuc.userRepo.GetByEmail(email)
	if err != nil {
		uuc.log.Error(context.Background(), "failed to get user by email", map[string]interface{}{
			"user_email": user.Email,
			"error":      err,
		})
		return pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get user by email",
		}
	}

	if user.VerificationCode != code {
		uuc.log.Info(context.Background(), "invalid verification code", map[string]interface{}{
			"user_id": user.ID,
			"code":    400,
		})
		return pkg.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid verification code",
		}
	}

	user.IsVerified = true
	user.VerificationCode = ""

	if err := uuc.userRepo.Update(&user); err != nil {
		uuc.log.Error(context.Background(), "failed to update user", map[string]interface{}{
			"user_id": user.ID,
			"error":   err,
		})
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

func (uuc *userUseCase) Auth(username, password string) (*domain.User, pkg.Response) {
	user, err := uuc.userRepo.GetByUsername(username)
	if err != nil {
		uuc.log.Error(context.Background(), "failed to get user by username", map[string]interface{}{
			"username": username,
			"error":    err,
		})
		return nil, pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get user by username",
		}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		uuc.log.Error(context.Background(), "failed to compare hash and password", map[string]interface{}{
			"user_id": user.ID,
			"error":   err,
		})
		return nil, pkg.Response{
			Code:    http.StatusUnauthorized,
			Message: "invalid password",
		}
	}

	return &user, pkg.Response{
		Code: http.StatusOK,
	}
}

func (uuc *userUseCase) ChangePass(userID string, newPassword string) pkg.Response {
	user, err := uuc.userRepo.GetByID(userID)
	if err != nil {
		uuc.log.Error(context.Background(), "failed to get user by id", map[string]interface{}{
			"user_id": userID,
			"error":   err,
		})
		return pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get user by username",
		}
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		uuc.log.Error(context.Background(), "failed to hash password", map[string]interface{}{
			"error": err,
		})
		return pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to hash password",
		}
	}

	user.Password = string(hashPass)

	if err := uuc.userRepo.Update(&user); err != nil {
		uuc.log.Error(context.Background(), "failed to update user", map[string]interface{}{
			"user_id": user.ID,
			"error":   err,
		})
		return pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update user",
		}
	}

	if err := uuc.sessionRepo.DeleteAllSessions(user.ID); err != nil {
		uuc.log.Error(context.Background(), "failed to delete sessions", map[string]interface{}{
			"user_id": user.ID,
			"error":   err,
		})
		return pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "failed to delete sessions by id",
		}
	}

	return pkg.Response{
		Code:    http.StatusOK,
		Message: "Password changed successfully",
	}
}
