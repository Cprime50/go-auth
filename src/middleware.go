package src

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if !strings.HasPrefix(tokenString, "Bearer ") {
			fmt.Printf("missing or wrong header %s", tokenString)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User Not logged in"})
			c.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := verifyToken(tokenString)
		if err != nil {
			fmt.Print("Invalid token", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User Not logged in"})
			c.Abort()
			return
		}

		username, ok := claims["user_name"].(string)
		if !ok {
			fmt.Print("Invalid token claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User Not logged in"})
			c.Abort()
			return
		}

		c.Set("user_name", username)
		c.Next()
	}
}
