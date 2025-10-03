// This file provides utility functions for password encryption and comparison.
package utils

// "golang.org/x/crypto/bcrypt" provides functions for hashing and comparing passwords using the bcrypt algorithm.
import (
	"golang.org/x/crypto/bcrypt"
)

// EncryptPassword hashes a password using the bcrypt algorithm.
// It takes a plain-text password as input and returns the hashed password and an error.
//
// @param password string - The plain-text password to be hashed.
// @return string - The hashed password.
// @return error - An error if one occurred during the hashing process.
func EncryptPassword(password string) (string, error) {
	// encryptedPassword is the hashed password.
	// bcrypt.GenerateFromPassword() hashes the password with a cost of 10.
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	// This checks if an error occurred while hashing the password.
	if err != nil {
		// If an error occurs, return an empty string and the error.
		return "", err
	}

	// The hashed password is converted to a string and returned.
	return string(encryptedPassword), nil
}

// CompareEncryptedPassword compares a hashed password with a plain-text password.
// It takes a hashed password and a plain-text password as input and returns a boolean indicating whether they match.
//
// @param encryptedPassword string - The hashed password.
// @param password string - The plain-text password.
// @return bool - True if the passwords match, false otherwise.
func CompareEncryptedPassword(encryptedPassword, password string) bool {
	// err is the result of comparing the hashed password with the plain-text password.
	// bcrypt.CompareHashAndPassword() compares a hashed password with its possible plaintext equivalent.
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	// The function returns true if the error is nil, indicating that the passwords match.
	return err == nil
}