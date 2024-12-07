package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	utils "github.com/OxytocinGroup/theca-backend/internal/api/utils/email"
	"github.com/OxytocinGroup/theca-backend/internal/config"
	"github.com/OxytocinGroup/theca-backend/internal/domain"
	services "github.com/OxytocinGroup/theca-backend/internal/usecase/interface"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userUseCase services.UserUseCase
}

// NewUserHandler returns a new UserHandler instance.
func NewUserHandler(usecase services.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: usecase,
	}
}

// @Register GoDoc
// @Summary Register a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body domain.UserRequest true "User"
// @Success 201 {object} domain.SuccessResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 409 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /register [post]
func (cr *UserHandler) Register(c *gin.Context) {
	var userRequest domain.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
		return
	}
	var user domain.User
	user.Email = userRequest.Email
	user.Password = userRequest.Password
	user.Username = userRequest.Username

	emailExists, err := cr.userUseCase.EmailExists(c, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check email existence",
		})
		return
	}
	if emailExists {
		c.JSON(http.StatusConflict, domain.ErrorResponse{
			Code:    http.StatusConflict,
			Message: "Email already exists",
		})
		return
	}

	usernameExists, err := cr.userUseCase.UsernameExists(c, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check username existence",
		})
		return
	}
	if usernameExists {
		c.JSON(http.StatusConflict, domain.ErrorResponse{
			Code:    http.StatusConflict,
			Message: "Username already exists",
		})
		return
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to hash password",
		})
		return
	}

	user.Password = string(hashPass)

	user.VerificationCode = strconv.Itoa(rand.Intn(900000) + 100000)
	user.IsVerified = false

	if err := cr.userUseCase.Create(c, &user); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
		})
		return
	}

	mail := utils.Mail{Email: user.Email, Code: user.VerificationCode, Username: user.Username}
	config, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to load config " + err.Error(),
		})
		return
	}
	if err := mail.SendVerificationEmail(config, user.Email, user.VerificationCode, user.Username); err != nil {
		fmt.Println(err)
		// #TODO LOGGER
		return
	}

	c.JSON(http.StatusCreated, domain.SuccessResponse{
		Message: "User registered successfully",
	})
}
