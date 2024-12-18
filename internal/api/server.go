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

	engine.POST("/register", userHandler.Register)
	engine.POST("/verify-email", userHandler.VerifyEmail)
	engine.POST("/login", userHandler.Login)
	engine.POST("/password-reset/request", userHandler.RequestPasswordReset)
	engine.POST("/password-reset/reset", userHandler.ResetPassword)

	// Auth middleware
	api := engine.Group("/api", middleware.AuthMiddleware(userHandler.SessionUseCase))
	// api.POST("/change-pass", userHandler.ChangePass)

	api.POST("/create-bookmark", bookmarkHandler.CreateBookmark)
	api.DELETE("/logout", userHandler.Logout)
	api.GET("/bookmarks", bookmarkHandler.GetBookmarks)
	api.DELETE("/bookmarks", bookmarkHandler.DeleteBookmark)
	api.POST("/bookmarks", bookmarkHandler.UpdateBookmark)
	api.GET("/user/verification-status", userHandler.CheckVerificationStatus)
	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	sh.engine.Run(":3000")
}
