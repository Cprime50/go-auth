package src

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserDetails(c *gin.Context) {
	username := c.GetString("user_name")
	fmt.Print("userid:", username)
	user, exists := users[username]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userDetails := map[string]interface{}{
		"username":    user.Username,
		"repositorys": getUserRepositories(username),
	}

	c.JSON(http.StatusOK, userDetails)
}
