package utils

import (
	// Import the bcrypt package from golang.org/x/crypto/bcrypt for secure password hashing.
	// bcrypt is a password hashing function designed to be computationally intensive,
	// making brute-force attacks more difficult.
	"golang.org/x/crypto/bcrypt"
)

// EncryptPassword takes a plain-text string (typically a password) and hashes it using the bcrypt algorithm.
// This function is crucial for securely storing sensitive data like passwords, as it prevents
// direct storage of the plain text, even if the database is compromised.
//
// Parameters:
// - password: The string to be encrypted (e.g., a user's password).
//
// Returns:
// - A string representing the bcrypt hash of the input data.
// - An error if the hashing process fails.
func EncryptPassword(password string) (string, error) {
	// GenerateFromPassword hashes the password using a cost factor of 10.
	// The cost factor determines how computationally expensive the hashing process is;
	// a higher cost factor makes it harder for attackers to crack hashes.
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		// If an error occurs during hashing, return an empty string and the error.
		return "", err
	}

	// Convert the byte slice hash to a string and return it.
	return string(encryptedPassword), nil
}

// CompareEncryptedPassword compares a plain-text password with a bcrypt-hashed password.
// This function is used during user authentication to verify if the provided password
// matches the stored hash without ever needing to decrypt the hash.
//
// Parameters:
// - encryptedPassword: The bcrypt hash retrieved from storage (e.g., from a database).
// - password: The plain-text password provided by the user during login.
//
// Returns:
// - true if the plain-text password matches the hashed password, indicating successful authentication.
// - false if they do not match or if an error occurs during the comparison.
func CompareEncryptedPassword(encryptedPassword, password string) bool {
	// CompareHashAndPassword compares a bcrypt hash with a plain-text password.
	// It returns nil if the password and hash match, and an error otherwise.
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	// Return true if err is nil (meaning the passwords match), otherwise return false.
	return err == nil
}
