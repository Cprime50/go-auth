package src

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var verifyCallbackTokens = make(map[string]VerifyCallback)

// Oauth2 config
var githubOAuthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
	ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	Endpoint:     github.Endpoint,
	RedirectURL:  "http://localhost:8000/oauth-callback/github",
}

func getOAuthConfig() *oauth2.Config {
	return githubOAuthConfig
}

func getUserInfo(accessToken string) (*User, error) {
	url := "https://api.github.com/user"
	userInfo, err := httpCall("GET", url, accessToken)
	if err != nil {
		return nil, fmt.Errorf("httpCall: %w", err)
	}

	username, ok := userInfo["name"].(string)
	if !ok {
		username = ""
	}
	return &User{
		Username: username,
		Password: "",
	}, nil
}

func httpCall(method, url, accessToken string) (map[string]interface{}, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %w", err)
	}
	defer resp.Body.Close()
	var userInfo map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return nil, fmt.Errorf("json.NewDecoder: %w", err)
	}
	return userInfo, nil
}

func generateState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func GithubLogin(c *gin.Context) {
	state, err := generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate state"})
		return
	}

	verifier := oauth2.GenerateVerifier()

	verifyCallbackTokens[state] = VerifyCallback{
		State:    state,
		Verifier: verifier,
		Expires:  time.Now().Add(10 * time.Minute),
	}

	url := githubOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	fmt.Print("Redirect URL:", url)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func getGithubOauthToken(ctx context.Context, code string, verifier string) (*oauth2.Token, error) {
	config := getOAuthConfig()
	oauthToken, err := config.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return nil, fmt.Errorf("error exchanging code for token: %w", err)
	}
	return oauthToken, nil
}

func OAuthCallBack(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	tokenData, exists := verifyCallbackTokens[state]
	if (!exists) || time.Now().After(tokenData.Expires) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired state"})
		return
	}

	delete(verifyCallbackTokens, state)

	token, err := getGithubOauthToken(c.Request.Context(), code, tokenData.Verifier)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	githubUser, err := getUserInfo(token.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, exists := users[githubUser.Username]
	if !exists {
		user = User{
			Username: githubUser.Username,
			Password: "",
		}
		users[user.Username] = user
	}

	accessToken, err := generateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	fmt.Print("access_token:", accessToken)
	c.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"access_token": accessToken,
		"user":         user,
	})
}
