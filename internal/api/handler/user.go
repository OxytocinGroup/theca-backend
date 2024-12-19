package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/OxytocinGroup/theca-backend/internal/usecase"
	"github.com/OxytocinGroup/theca-backend/pkg"
	"github.com/OxytocinGroup/theca-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	UserUseCase    usecase.UserUseCase
	SessionUseCase usecase.SessionUseCase
	Logger         logger.Logger
}

func NewUserHandler(usecase usecase.UserUseCase, sessionUseCase usecase.SessionUseCase, log logger.Logger) *UserHandler {
	return &UserHandler{
		UserUseCase:    usecase,
		SessionUseCase: sessionUseCase,
		Logger:         log,
	}
}

// @Register GoDoc
// @Summary Register a new user
// @Description Register a new user
// @Tags User
// @Accept json
// @Produce json
// @Param request body struct{email string, username string, password string} true "Username, email and password"
// @Success 201 {object} pkg.Response
// @Failure 409 {object} pkg.Response
// @Failure 500 {object} pkg.Response
// @Failure 400 {object} pkg.Response
// @Router /register [post]
func (uh *UserHandler) Register(c *gin.Context) {
	var userRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Username string `json:"username" binding:"required,min=3"`
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, pkg.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
		return
	}
	resp := uh.UserUseCase.Register(userRequest.Email, userRequest.Password, userRequest.Username)
	c.JSON(resp.Code, resp)
}

// @VerifyEmail GoDoc
// @Summary Verify email address
// @Description This endpoint allows a user to verify their email address by providing the email and verification code.
// @Tags User
// @Accept json
// @Produce json
// @Param request body struct{email string, code string} true "email and verification code"
// @Success 200 {object} pkg.Response "Email verified successfully"
// @Failure 400 {object} pkg.Response "Bad request - Invalid input"
// @Failure 404 {object} pkg.Response "Email or verification code not found"
// @Failure 500 {object} pkg.Response "Internal server error"
// @Router /verify-email [post]
func (uh *UserHandler) VerifyEmail(c *gin.Context) {
	var verifyReq struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&verifyReq); err != nil {
		c.JSON(http.StatusBadRequest, pkg.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
		return
	}

	resp := uh.UserUseCase.VerifyEmail(verifyReq.Email, verifyReq.Code)
	c.JSON(resp.Code, resp)
}

// @Login GoDoc
// @Summary User login
// @Description This endpoint allows a user to log in using their username and password. If already logged in, a conflict response is returned.
// @Tags User
// @Accept  json
// @Produce  json
// @Param request body struct{username string, password string} true "Username and password"
// @Success 200 {object} pkg.Response "Login successful"
// @Failure 400 {object} pkg.Response "Bad request - Invalid input"
// @Failure 401 {object} pkg.Response "Unauthorized - Invalid username or password"
// @Failure 409 {object} pkg.Response "Conflict - User already logged in"
// @Failure 500 {object} pkg.Response "Internal server error"
// @Router /login [post]
func (uh *UserHandler) Login(c *gin.Context) {
	session, err := c.Cookie("session_id")
	if err == nil {
		userID, err := uh.SessionUseCase.ValidateSession(session)
		if err == nil {
			uh.Logger.Info(context.Background(), "user tryed to login when already logged", map[string]any{"user_id": userID})
			c.JSON(http.StatusConflict, pkg.Response{
				Code:    http.StatusConflict,
				Message: "user already logged in",
			})
			return
		}
	}

	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, pkg.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
		return
	}

	user, resp := uh.UserUseCase.Auth(req.Username, req.Password)
	if resp.Code != 200 {
		c.JSON(resp.Code, resp)
		return
	}

	sessionID := uuid.New().String()
	if err := uh.SessionUseCase.CreateSession(sessionID, user.ID, time.Now().Add(24*time.Hour)); err != nil {
		uh.Logger.Error(context.Background(), "failed to create session", map[string]any{
			"error": err,
		})
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
// @Summary User logout
// @Description This endpoint allows a user to log out by deleting all active sessions associated with the user.
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} pkg.Response "Logout successful"
// @Failure 500 {object} pkg.Response "Internal server error"
// @Router /logout [post]
func (uh *UserHandler) Logout(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := uh.SessionUseCase.DeleteAllSessions(userID); err != nil {
		uh.Logger.Error(context.Background(), "failed to delete sessions", map[string]any{
			"user_id": userID,
			"error":   err,
		})
		c.JSON(http.StatusInternalServerError, pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "failed to delete sessions",
		})
		return
	}

	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, pkg.Response{
		Code:    http.StatusOK,
		Message: "Logout successful",
	})
}

// @ChangePass GoDoc
// @Summary Change user password
// @Description Changes the user's password and deleting all sessions for this user
// @Tags User
// @Produce json
// @Param request body struct{password string} true "New password"
// @Success 200 {object} pkg.Response
// @Failure 400 {object} pkg.Response
// @Failure 401 {object} pkg.Response
// @Failure 500 {object} pkg.Response
// @Router /api/change-pass [post]
// @Security CookieAuth
func (uh *UserHandler) ChangePass(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, pkg.Response{
			Code: http.StatusUnauthorized,
		})
		return
	}

	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, pkg.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
		return
	}

	resp := uh.UserUseCase.ChangePass(fmt.Sprint(userID), req.Password)
	c.JSON(resp.Code, resp)
}

// CheckVerificationStatus godoc
// @Summary Check verification status of the user
// @Description Check if the user's verification status is complete
// @Tags User
// @Accept json
// @Produce json
// @Security CookieAuth
// @Success 204 {object} pkg.Response "User verification status checked successfully, no content"
// @Failure 500 {object} pkg.Response "Internal server error"
// @Router /api/user/verification-status [get]
func (uh *UserHandler) CheckVerificationStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	exists, resp := uh.UserUseCase.CheckVerificationStatus(userID)
	if resp.Code != 200 {
		c.JSON(resp.Code, resp)
		return
	}
	c.JSON(resp.Code, map[string]bool{
		"status": exists,
	})
}

// @RequestPasswordReset godoc
// @Summary Request a password reset
// @Description This endpoint allows a user to request a password reset by providing their email.
// @Tags User
// @Accept  json
// @Produce  json
// @Param request body struct{email string} true "User email"
// @Success 200 {object} pkg.Response "Password reset email sent"
// @Failure 400 {object} pkg.Response "Bad request - Invalid input"
// @Failure 404 {object} pkg.Response "Email not found"
// @Failure 500 {object} pkg.Response "Internal server error"
// @Router /password-reset/request [post]
func (uh *UserHandler) RequestPasswordReset(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, pkg.Response{Code: 400, Message: "bad request"})
		return
	}

	resp := uh.UserUseCase.GetResetPassword(req.Email)
	c.JSON(resp.Code, resp)
}

// ResetPassword godoc
// @Summary Reset user password
// @Description This endpoint allows a user to reset their password by providing a valid reset token and a new password.
// @Tags User
// @Accept  json
// @Produce  json
// @Param request body struct{token string; password string} true "Reset token and new password"
// @Success 200 {object} pkg.Response "Password reset successfully"
// @Failure 400 {object} pkg.Response "Bad request - Invalid input"
// @Failure 404 {object} pkg.Response "Reset token not found or expired"
// @Failure 500 {object} pkg.Response "Internal server error"
// @Router /password-reset/reset [post]
func (uh *UserHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, pkg.Response{Code: 400, Message: "bad request"})
		return
	}

	resp := uh.UserUseCase.ResetPassword(req.Token, req.Password)
	c.JSON(resp.Code, resp)
}
