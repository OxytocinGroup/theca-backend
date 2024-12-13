package handler

import (
	"net/http"
	"time"

	"github.com/OxytocinGroup/theca-backend/internal/usecase"
	"github.com/OxytocinGroup/theca-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	UserUseCase    usecase.UserUseCase
	SessionUseCase usecase.SessionUseCase
}

func NewUserHandler(usecase usecase.UserUseCase, sessionUseCase usecase.SessionUseCase) *UserHandler {
	return &UserHandler{
		UserUseCase:    usecase,
		SessionUseCase: sessionUseCase,
	}
}

// @Register GoDoc
// @Summary Register a new user
// @Description Register a new user
// @Tags User
// @Accept json
// @Produce json
// @Param user body pkg.UserRequest true "User"
// @Success 201 {object} pkg.Response
// @Failure 409 {object} pkg.Response
// @Failure 500 {object} pkg.Response
// @Failure 400 {object} pkg.Response
// @Router /register [post]
// @Security ApiKeyAuth
func (cr *UserHandler) Register(c *gin.Context) {
	var userRequest pkg.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, pkg.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
		return
	}
	resp := cr.UserUseCase.Register(userRequest.Email, userRequest.Password, userRequest.Username)
	c.JSON(resp.Code, resp)
}

// @VerifyEmail GoDoc
// @Summary Verify email
// @Description Verify email
// @Tags User
// @Accept json
// @Produce json
// @Param verifyReq body pkg.VerifyRequest true "VerifyRequest"
// @Success 200 {object} pkg.Response
// @Failure 400 {object} pkg.Response
// @Failure 500 {object} pkg.Response
// @Router /api/verify-email [post]
// @Security ApiKeyAuth
func (cr *UserHandler) VerifyEmail(c *gin.Context) {
	var verifyReq pkg.VerifyRequest
	if err := c.ShouldBindJSON(&verifyReq); err != nil {
		c.JSON(http.StatusBadRequest, pkg.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
		return
	}

	resp := cr.UserUseCase.VerifyEmail(verifyReq.Email, verifyReq.Code)
	c.JSON(resp.Code, resp)
}

// @Login GoDoc
// @Summary User login
// @Description Authenticates a user and initiates a session by setting a session cookie.
// @Tags User
// @Accept json
// @Produce json
// @Param user body pkg.LoginRequest true "User"
// @Success 200 {object} pkg.Response
// @Failure 400 {object} pkg.Response
// @Failure 401 {object} pkg.Response
// @Failure 500 {object} pkg.Response
// @Router /api/login [post]
// @Security ApiKeyAuth
func (cr *UserHandler) Login(c *gin.Context) {
	var req pkg.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, pkg.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
		return
	}

	user, err := cr.UserUseCase.Auth(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, pkg.Response{
			Code:    http.StatusUnauthorized,
			Message: "Invalid credentials",
		})
		return
	}

	sessionID := uuid.New().String()
	if err := cr.SessionUseCase.CreateSession(sessionID, user.ID, time.Now().Add(24*time.Hour)); err != nil {
		c.JSON(http.StatusInternalServerError, pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create session" + err.Error(),
		})
		return
	}

	c.SetCookie("session_id", sessionID, 3600*24, "/", "", false, true)
	c.JSON(http.StatusOK, pkg.Response{
		Code:    http.StatusOK,
		Message: "Login successful",
	})
}

// @Logout GoDoc
// @Summary Logout a user
// @Description Logout a user by deleting the session and removing the session cookie.
// @Tags User
// @Produce json
// @Success 200 {object} pkg.Response
// @Failure 401 {object} pkg.Response
// @Failure 500 {object} pkg.Response
// @Router /api/logout [delete]
// @Security ApiKeyAuth
func (cr *UserHandler) Logout(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, pkg.Response{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		})
		return
	}

	if err := cr.SessionUseCase.DeleteSession(sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete session",
		})
		return
	}

	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, pkg.Response{
		Code:    http.StatusOK,
		Message: "Logout successful",
	})

}
