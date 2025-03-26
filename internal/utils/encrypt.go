package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

func SHA1Hash(input string) (string, error) {
	hash := sha1.New()

	_, err := hash.Write([]byte(input))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
