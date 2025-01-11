package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/OxytocinGroup/theca-backend/internal/usecase"
	"github.com/OxytocinGroup/theca-backend/pkg"
	"github.com/OxytocinGroup/theca-backend/pkg/cerr"
	"github.com/OxytocinGroup/theca-backend/pkg/logger"
	"github.com/OxytocinGroup/theca-backend/pkg/requests"
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
// @Param request body requests.RegisterRequest true "Username, email and password"
// @Success 201 {object} pkg.Response
// @Failure 409 {object} pkg.Response
// @Failure 500 {object} pkg.Response
// @Failure 400 {object} pkg.Response
// @Router /user/register [post]
func (uh *UserHandler) Register(c *gin.Context) {
	var req requests.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uh.Logger.Info(context.Background(), "Register: bad request", map[string]any{"error": err})
		c.JSON(http.StatusBadRequest, pkg.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   cerr.ErrInvalidBody,
		})
		return
	}
	resp := uh.UserUseCase.Register(req.Email, req.Password, req.Username)
	c.JSON(resp.Code, resp)
}

// @VerifyEmail GoDoc
// @Summary Verify email address
// @Description This endpoint allows a user to verify their email address by providing the email and verification code.
// @Tags User
// @Accept json
// @Produce json
// @Param request body requests.EmailVerifyRequest true "verification code"
// @Success 200 {object} pkg.Response "Email verified successfully"
// @Failure 400 {object} pkg.Response "Bad request - Invalid input"
// @Failure 404 {object} pkg.Response "Verification code not found"
// @Failure 500 {object} pkg.Response "Internal server error"
// @Router /user/verify-email [post]
func (uh *UserHandler) VerifyEmail(c *gin.Context) {
	var req requests.EmailVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uh.Logger.Info(context.Background(), "Verify email: bad request", map[string]any{"error": err})
		c.JSON(http.StatusBadRequest, pkg.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   cerr.ErrInvalidBody,
		})
		return
	}

	resp := uh.UserUseCase.VerifyEmail(req.Code)
	c.JSON(resp.Code, resp)
}

// @Login GoDoc
// @Summary User login
// @Description This endpoint allows a user to log in using their username and password. If already logged in, a conflict response is returned.
// @Tags User
// @Accept  json
// @Produce  json
// @Param request body requests.LoginRequest true "Username and password"
// @Success 200 {object} pkg.LoginResponse "Login successful"
// @Failure 400 {object} pkg.Response "Bad request - Invalid input"
// @Failure 401 {object} pkg.Response "Unauthorized - Invalid username or password or not verified"
// @Failure 409 {object} pkg.Response "Conflict - User already logged in"
// @Failure 500 {object} pkg.Response "Internal server error"
// @Router /user/login [post]
func (uh *UserHandler) Login(c *gin.Context) {
	session, err := c.Cookie("session_id")
	if err == nil {
		userID, err := uh.SessionUseCase.ValidateSession(session)
		if err == nil {
			uh.Logger.Info(context.Background(), "Login: user tryed to login when already logged", map[string]any{"user_id": userID})
			c.JSON(http.StatusConflict, pkg.Response{
				Code:    http.StatusConflict,
				Message: "User already logged in",
				Error:   cerr.ErrUserLogined,
			})
			return
		}
	}

	var req requests.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		uh.Logger.Info(context.Background(), "Login: bad request", map[string]any{"error": err})
		c.JSON(http.StatusBadRequest, pkg.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   cerr.ErrInvalidBody,
		})
		return
	}

	user, resp := uh.UserUseCase.Auth(req.Username, req.Password)
	if resp.Code != 200 {
		c.JSON(resp.Code, resp)
		return
	}

	verified, resp := uh.UserUseCase.CheckVerificationStatus(req.Username)
	if resp.Code != http.StatusOK {
		c.JSON(http.StatusUnauthorized, resp)
		return
	}
	if !verified {
		uh.Logger.Info(context.Background(), "Login: email is not verified", map[string]any{"user_id": req.Username})
		c.JSON(http.StatusUnauthorized, pkg.Response{Code: http.StatusUnauthorized, Message: "not verified", Error: cerr.ErrEmailNotVerified})
		return
	}

	sessionID := uuid.New().String()
	if err := uh.SessionUseCase.CreateSession(sessionID, user.ID, time.Now().Add(24*time.Hour)); err != nil {
		uh.Logger.Error(context.Background(), "Login: failed to create session", map[string]any{
			"error": err,
		})
		c.JSON(http.StatusInternalServerError, pkg.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create session" + err.Error(),
		})
		return
	}

	c.SetCookie("session_id", sessionID, 3600*24, "/", "", true, true)
	c.JSON(http.StatusOK, pkg.LoginResponse{
		Code:     http.StatusOK,
		Message:  "Login successful",
		Username: user.Username,
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
// @Router /api/user/logout [delete]
func (uh *UserHandler) Logout(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := uh.SessionUseCase.DeleteAllSessions(userID); err != nil {
		uh.Logger.Error(context.Background(), "Logut: failed to delete sessions", map[string]any{
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

// func (uh *UserHandler) ChangePass(c *gin.Context) {
// 	userID, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, pkg.Response{
// 			Code: http.StatusUnauthorized,
// 		})
// 		return
// 	}

// 	var req requests.ChangePasswordRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, pkg.Response{
// 			Code:    http.StatusBadRequest,
// 			Message: "Invalid request body",
// 		})
// 		return
// 	}

// 	resp := uh.UserUseCase.ChangePass(fmt.Sprint(userID), req.Password)
// 	c.JSON(resp.Code, resp)
// }

// @RequestPasswordReset godoc
// @Summary Request a password reset
// @Description This endpoint allows a user to request a password reset by providing their email.
// @Tags User
// @Accept  json
// @Produce  json
// @Param request body requests.RequestPasswordReset true "User email"
// @Success 200 {object} pkg.Response "Password reset email sent"
// @Failure 400 {object} pkg.Response "Bad request - Invalid input"
// @Failure 404 {object} pkg.Response "Email not found"
// @Failure 500 {object} pkg.Response "Internal server error"
// @Router /user/password-reset/request [post]
func (uh *UserHandler) RequestPasswordReset(c *gin.Context) {
	var req requests.RequestPasswordReset
	if err := c.ShouldBindJSON(&req); err != nil {
		uh.Logger.Info(context.Background(), "Request password reset: bad request", map[string]any{"error": err})
		c.JSON(http.StatusBadRequest, pkg.Response{Code: http.StatusBadRequest, Message: "bad request", Error: cerr.ErrInvalidBody})
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
// @Param request body requests.ResetPassword true "Reset token and new password"
// @Success 200 {object} pkg.Response "Password reset successfully"
// @Failure 400 {object} pkg.Response "Bad request - Invalid input"
// @Failure 404 {object} pkg.Response "Reset token not found or expired"
// @Failure 500 {object} pkg.Response "Internal server error"
// @Router /user/password-reset/reset [post]
func (uh *UserHandler) ResetPassword(c *gin.Context) {
	var req requests.ResetPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		uh.Logger.Info(context.Background(), "Reset password: bad request", map[string]any{"error": err})
		c.JSON(400, pkg.Response{Code: 400, Message: "Bad request", Error: cerr.ErrInvalidBody})
		return
	}

	resp := uh.UserUseCase.ResetPassword(req.Token, req.Password)
	c.JSON(resp.Code, resp)
}

// RequestVerificationToken godoc
// @Summary Resend Verification Token
// @Description Resends a verification token for the provided username.
// @Tags User
// @Accept json
// @Produce json
// @Param request body requests.RequestVerificationToken true "Request body containing the username"
// @Success 200 {object} pkg.Response "Verification token resent successfully"
// @Failure 400 {object} pkg.Response "Bad request, invalid input"
// @Failure 403 {object} pkg.Response "User has verified email"
// @Failure 404 {object} pkg.Response "User not found"
// @Failure 500 {object} pkg.Response "Internal server error"
// @Router /user/verify-email/request [post]
func (uh *UserHandler) RequestVerificationToken(c *gin.Context) {
	var req requests.RequestVerificationToken
	if err := c.ShouldBindJSON(&req); err != nil {
		uh.Logger.Info(context.Background(), "Resend token: bad request", map[string]any{"error": err})
		c.JSON(http.StatusBadRequest, pkg.Response{Code: http.StatusBadRequest, Message: "Bad request", Error: cerr.ErrInvalidBody})
		return
	}

	resp := uh.UserUseCase.ResendVerificationToken(req.Username)
	c.JSON(resp.Code, resp)
}
