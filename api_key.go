package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
)

func GenerateAPIKey(length int, prefix string) (string, error) {
	if length < 32 {
		return "", errors.New("minimum 32-character of api key")
	}

	byteLength := length / 2

	apiKeyBytes := make([]byte, byteLength)

	_, err := rand.Read(apiKeyBytes)
	if err != nil {
		return "", err
	}

	apiKey := hex.EncodeToString(apiKeyBytes)

	if prefix != "" {
		apiKey = fmt.Sprintf("%s%s_", prefix, apiKey)
	}

	return apiKey, nil
}
