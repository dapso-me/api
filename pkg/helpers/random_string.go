package helpers

import (
	"crypto/rand"
	"fmt"
)

func GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("helper: generate random string: length must be greater than zero")
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("helper: generate random string: %w", err)
	}

	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}

	return string(bytes), nil
}
