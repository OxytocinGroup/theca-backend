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

func NewServerHTTP(userHandler *handler.UserHandler) *ServerHTTP {
	engine := gin.New()

	engine.Use(gin.Logger())

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	engine.POST("/register", userHandler.Register)
	engine.POST("/verify-email", userHandler.VerifyEmail)
	engine.POST("/login", userHandler.Login)
	engine.DELETE("/logout", userHandler.Logout)
	// Auth middleware
	api := engine.Group("/api", middleware.AuthMiddleware(userHandler.SessionUseCase))
	api.GET("/protected", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			return
		}
		c.JSON(200, gin.H{"message": "Welcome to your profile!", "user_id": userID})
	})
	// api.POST("/change-pass", userHandler.ChangePass)
	
	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	sh.engine.Run(":3000")
}
