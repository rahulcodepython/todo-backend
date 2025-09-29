package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rahulcodepython/todo-backend/backend/config"
)

func createToken(userId string, cfg *config.Config) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": userId,
			"exp":    time.Now().Add(cfg.JWT.Expires).Unix(),
		})

	tokenString, err := token.SignedString(cfg.JWT.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
