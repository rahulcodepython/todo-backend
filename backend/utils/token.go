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

func CreateToken(userId string, cfg *config.Config) *Token {
	token := Token{
		Token:     "",
		ExpiresAt: time.Now().Add(cfg.JWT.Expires),
	}

	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(cfg.JWT.Expires).Unix(),
		"iat":     time.Now().Unix(),
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenClaims.SignedString([]byte(cfg.JWT.SecretKey))
	if err != nil {
		return nil
	}

	token.Token = tokenString
	return &token
}
