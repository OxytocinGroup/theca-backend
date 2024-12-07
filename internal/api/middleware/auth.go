package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// AuthorizationMiddleware is a middleware function for Gin that
// checks the Authorization header for a Bearer token. It validates
// the token, and if the validation fails, it aborts the request
// with a 401 Unauthorized status. This middleware should be used
// to protect routes that require authentication.
func AuthorizationMiddleware(c *gin.Context) {
	s := c.Request.Header.Get("Authorization")

	token := strings.TrimPrefix(s, "Bearer ")

	if err := validateToken(token); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

// validateToken validates the given JWT token. The token is verified
// against a secret "secret". If the token is invalid, an error is
// returned.
func validateToken(token string) error {
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte("secret"), nil
	})

	return err
}

// LoginHandler is a Gin handler that accepts a JSON payload with a "user"
// and "pass" key. If the credentials are invalid, it aborts the request with a
// 401 Unauthorized status. If the credentials are valid, it generates a JWT
// token and returns it in the response body.
func LoginHandler(c *gin.Context) {
	// implement login logic here
	// user := c.PostForm("user")
	// pass := c.PostForm("pass")

	var json struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	}

	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	user := json.User
	pass := json.Pass

	// Throws Unauthorized error
	if user != "john" || pass != "lark" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"name":  "John Lark",
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
	// 	ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
	// })

	ss, err := token.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"token": ss,
	})
}
