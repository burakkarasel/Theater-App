package util

import "golang.org/x/crypto/bcrypt"

// HashPassword takes a password string as input an returns the hashed password and possible error
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CompareHashedPassword takes password and hashed password as input and returns a possible error
func CompareHashedPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
