package main

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// Hash the password with a cost factor of bcrypt.DefaultCost (10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, password string) bool {
	// Compare the hashed password with the plain text password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

