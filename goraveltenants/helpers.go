package goraveltenants

import (
	"crypto/rand"
	"encoding/base64"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func GenerateSecureRandomString(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(randomBytes)[:length], nil
}

func MergeConfigMaps(first map[string]interface{}, second map[string]interface{}) map[string] interface{} {
	for key, val := range second {
		first[key] = val
	}

	return first
}
