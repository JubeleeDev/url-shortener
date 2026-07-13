package shortener

import (
	"fmt"
	"math/rand/v2"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func getAlphabetChar(index int) string {
	if index < 0 || index >= len(alphabet) {
		return ""
	}

	return string(alphabet[index])
}

func isValidCodeLength(length int) bool {
	return length > 0
}

func getRandomAlphabetChar() string {
	index := rand.IntN(len(alphabet))

	return getAlphabetChar(index)
}

func GenerateCode(length int) (string, error) {
	if !isValidCodeLength(length) {
		return "", fmt.Errorf("error: %w", ErrInvalidCodeLength)
	}

	result := ""

	for i := 0; i < length; i++ {
		result += getRandomAlphabetChar()
	}

	return result, nil
}
