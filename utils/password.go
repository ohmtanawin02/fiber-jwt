package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

func ComparePassword(oldPassword, newPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(oldPassword), []byte(newPassword))
	return err == nil
}
