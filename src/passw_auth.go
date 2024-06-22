package src

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var users = make(map[string]User)

func Register(c *gin.Context) {
	var signUpInfo struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&signUpInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if _, exists := users[signUpInfo.Username]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signUpInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	user := User{
		Username: signUpInfo.Username,
		Password: string(hashedPassword),
	}
	users[user.Username] = user

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}

func Login(c *gin.Context) {
	var loginInfo struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, exists := users[loginInfo.Username]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password))
	if err != nil {
		fmt.Println("Error comparing hashed password:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := generateToken(user.Username)
	if err != nil {
		fmt.Printf("Error creating tokens for user %s: %v\n", user.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not log in"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"access_token": token,
	})
}
