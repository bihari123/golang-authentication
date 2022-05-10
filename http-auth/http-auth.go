package httpauth

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func ConvertToBase64() {
	fmt.Println(base64.StdEncoding.EncodeToString([]byte("user:pass")))
}

func HashPassword(password string) (hashedPass []byte, err error) {
	hashedPass, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // more the cost parameter is , the more secure the hash is but it is also more time consuming
	if err != nil {
		return nil, fmt.Errorf("Error while generating bcrypt hash from password: %w", err)
	}
	return
}

// bcrypt has a function that returns the same time for both success and failure, so it is not easy to crack

func ComparePassword(password string, hashedPass []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPass,[]byte(password)) 
	if err != nil {
		return fmt.Errorf("Invalid password: %w", err)
	}
	return nil
}
