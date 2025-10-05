package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func IsAuthenticated(c *gin.Context) {
	var secret = []byte(os.Getenv("SECRET"))
	cookie, err := c.Cookie("jwt")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no cookie"})
		c.Abort()
		return
	}

	token, _ := jwt.Parse(cookie, func(t *jwt.Token) (any, error) {
		return secret, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}

	userID := uint(claims["userID"].(float64))
	c.Set("userID", userID)


	c.Next()
}
