package src

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	tokenSecretKey []byte
	tokenTTL       time.Duration
)

func LoadSecrets() {
	var err error
	tokenSecretKey = []byte(os.Getenv("ACCESS_SECRET_KEY"))
	tokenTTLStr := os.Getenv("TOKEN_TTL")

	if len(tokenSecretKey) == 0 || tokenTTLStr == "" {
		panic("SECRETs in environment variables are not set correctly")
	}

	tokenTTL, err = time.ParseDuration(tokenTTLStr)
	if err != nil {
		panic(fmt.Sprintf("Error parsing TOKEN_TTL: %v", err))
	}
}

func generateToken(username string) (string, error) {
	expTime := time.Now().Add(tokenTTL).Unix()
	claims := jwt.MapClaims{
		"user_name": username,
		"exp":       expTime,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tokenSecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func verifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return tokenSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("could not parse claims")
	}

	return claims, nil
}
