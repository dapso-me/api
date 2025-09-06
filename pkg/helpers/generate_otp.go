package helpers

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTP(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("helper: generate otp: length must be greater than zero")
	}

	digits := "0123456789"

	result := make([]byte, length)

	for i := range length {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("helper: generate otp: %w", err)
		}

		result[i] = digits[randomIndex.Int64()]
	}

	return string(result), nil
}
