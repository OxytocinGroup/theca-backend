package middleware

import (
	"net/http"

	"github.com/OxytocinGroup/theca-backend/internal/usecase"
	"github.com/OxytocinGroup/theca-backend/pkg/cerr"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(sessionUC usecase.SessionUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": cerr.ErrMissingCookie,})
			c.Abort()
			return
		}

		userID, err := sessionUC.ValidateSession(sessionID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "details": err.Error()})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
