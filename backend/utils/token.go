package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rahulcodepython/todo-backend/backend/config"
)

type Token struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func CreateToken(userId string, cfg *config.Config) (*Token, error) {
	token := Token{
		Token:     "",
		ExpiresAt: time.Now().Add(cfg.JWT.Expires),
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": userId,
			"exp":    time.Now().Add(cfg.JWT.Expires).Unix(),
		})

	tokenString, err := jwt.SignedString(cfg.JWT.SecretKey)
	if err != nil {
		return nil, err
	}

	token.Token = tokenString
	return &token, nil
}
