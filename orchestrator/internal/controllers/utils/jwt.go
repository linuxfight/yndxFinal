package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateToken(userID string, expires time.Time, secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": expires.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
