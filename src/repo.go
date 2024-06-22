package src

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var repos = make(map[string]Repo)

func CreateRepo(c *gin.Context) {
	var repoInfo struct {
		Name    string `json:"name" binding:"required"`
		Details string `json:"details" binding:"required"`
	}
	if err := c.ShouldBindJSON(&repoInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input, 'name' field is required"})
		return
	}

	username := c.GetString("user_name")
	repoID := uuid.New().String()
	repo := Repo{
		ID:       repoID,
		Name:     repoInfo.Name,
		Username: username,
	}
	repos[repo.ID] = repo

	c.JSON(http.StatusCreated, gin.H{
		"message": "Repository created successfully",
		"repo":    repo,
	})
}

func getUserRepositories(username string) []Repo {
	var userRepos []Repo
	for _, repo := range repos {
		if username == repo.Username {
			userRepos = append(userRepos, repo)
		}
	}
	return userRepos
}

func EditRepo(c *gin.Context) {
	repoId := c.Param("id")
	username := c.GetString("user_name")
	var repoInfo struct {
		Name    string `json:"name" binding:"required"`
		Details string `json:"details" binding:"required"`
	}
	if err := c.ShouldBindJSON(&repoInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input, 'name' field is required"})
		return
	}

	repo, exists := repos[repoId]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}
	if repo.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to edit this repository"})
		return
	}

	repo.Name = repoInfo.Name
	repo.Details = repoInfo.Details
	repos[repoId] = repo

	c.JSON(http.StatusOK, gin.H{
		"message": "Repository edited successfully",
		"repo":    repo,
	})
}

func DeleteRepo(c *gin.Context) {
	repoID := c.Param("id")
	username := c.GetString("user_name")
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	repo, exists := repos[repoID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	if repo.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this repository"})
		return
	}

	delete(repos, repoID)

	c.JSON(http.StatusOK, gin.H{"message": "Repository deleted successfully"})
}
