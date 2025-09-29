package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(data string) (string, error) {
	encryptedData, err := bcrypt.GenerateFromPassword([]byte(data), 10)
	if err != nil {
		return "", err
	}

	return string(encryptedData), nil
}

func CompareEncryptedData(encryptedData, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedData), []byte(password))
	return err == nil
}
