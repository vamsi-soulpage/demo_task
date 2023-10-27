package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
)

var jwtSecret = []byte("soulpage")

func ExtractUsernameFromToken(c *gin.Context) (string, error) {
	// Get the Authorization header from the request
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		return "", errors.New("Token is missing")
	}

	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || authParts[0] != "Bearer" {
		return "", errors.New("Invalid token format")
	}

	tokenString := authParts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("Invalid token")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", errors.New("Username not found in token claims")
	}

	return username, nil
}
