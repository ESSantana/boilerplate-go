package utils

import (
	"fmt"
	"os"
	"strings"
)

func RetrieveSecretValue(key string) string {
	if strings.HasSuffix(key, "_FILE") {
		secretPath := os.Getenv(key)
		data, err := os.ReadFile(secretPath)
		if err != nil {
			fmt.Println("error reading secret file: ", err.Error())
			return ""
		}
		return string(data)
	}
	return os.Getenv(key)
}
