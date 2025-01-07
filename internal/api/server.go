package http

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/OxytocinGroup/theca-backend/cmd/api/docs"
	handler "github.com/OxytocinGroup/theca-backend/internal/api/handler"
	"github.com/OxytocinGroup/theca-backend/internal/api/middleware"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(userHandler *handler.UserHandler, bookmarkHandler *handler.BookmarkHandler) *ServerHTTP {
	engine := gin.New()

	engine.Use(gin.Logger())

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	engine.POST("/user/register", userHandler.Register)
	engine.POST("/user/verify-email", userHandler.VerifyEmail)
	engine.POST("/user/login", userHandler.Login)
	engine.POST("/user/password-reset/request", userHandler.RequestPasswordReset)
	engine.POST("/user/password-reset/reset", userHandler.ResetPassword)

	// Auth middleware
	api := engine.Group("/api", middleware.AuthMiddleware(userHandler.SessionUseCase))
	// api.POST("/change-pass", userHandler.ChangePass)

	api.DELETE("/user/logout", userHandler.Logout)
	api.POST("/bookmarks/create", bookmarkHandler.CreateBookmark)
	api.GET("/bookmarks/get", bookmarkHandler.GetBookmarks)
	api.DELETE("/bookmarks/delete", bookmarkHandler.DeleteBookmark)
	api.POST("/bookmarks/update", bookmarkHandler.UpdateBookmark)
	// api.GET("/user/verification-status", userHandler.CheckVerificationStatus)
	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	sh.engine.Run(":3000")
}
