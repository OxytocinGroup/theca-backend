package handler

import (
	"net/http"

	services "github.com/OxytocinGroup/theca-backend/internal/usecase/interface"
	"github.com/OxytocinGroup/theca-backend/pkg"
	"github.com/gin-gonic/gin"
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
// @Description Register a new user
// @Tags User
// @Accept json
// @Produce json
// @Param user body pkg.UserRequest true "User"
// @Success 201 {object} pkg.Response
// @Failure 409 {object} pkg.Response
// @Failure 500 {object} pkg.Response
// @Failure 400 {object} pkg.Response
// @Router /api/users [post]
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

	resp := cr.userUseCase.Register(userRequest.Email, userRequest.Password, userRequest.Username)
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

	resp := cr.userUseCase.VerifyEmail(verifyReq.Email, verifyReq.Code)
	c.JSON(resp.Code, resp)
}
