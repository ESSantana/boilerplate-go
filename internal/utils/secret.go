package utils

import "os"

func RetrieveSecretValue(key string) string {
	secretPath := os.Getenv(key)
	if secretPath == "" {
		return ""
	}

	data, err := os.ReadFile(secretPath)
	if err != nil {
		return ""
	}

	return string(data)
}
