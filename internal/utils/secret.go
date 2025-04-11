package utils

import (
	"os"
	"strings"
)

func RetrieveSecretValue(key string) string {
	secretPath := os.Getenv(key)
	if secretPath == "" {
		return ""
	}

	if !strings.Contains(key, "/") {
		// is not a secret path, return the value directly
		return secretPath
	}

	data, err := os.ReadFile(secretPath)
	if err != nil {
		return ""
	}

	return string(data)
}
