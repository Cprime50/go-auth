package main

import (
	"github.com/Cprime50/go-auth/src"
	"github.com/gin-gonic/gin"
)

func main() {
	src.LoadSecrets()

	r := gin.Default()

	r.POST("/register", src.Register)
	r.POST("/login", src.Login)
	r.GET("/login/github", src.GithubLogin)
	r.GET("/oauth-callback/github", src.OAuthCallBack)

	// Authorized group
	auth := r.Group("/auth")
	auth.Use(src.AuthMiddleware())
	{
		auth.GET("/user", src.GetUserDetails)
		auth.POST("/repo", src.CreateRepo)
		auth.PUT("/repo/:id", src.EditRepo)
		auth.DELETE("/repo/:id", src.DeleteRepo)
	}

	r.Run(":8000")
}
