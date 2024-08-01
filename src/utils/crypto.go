package utils

import (
	"crypto/rand"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/hex"
	"strings"
)

const (
	possible    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	possibleLen = len(possible)
)

func generateRandomByteSlice(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateRandomId() (string, error) {
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateCodeVerifier(length int) (string, error) {
	text := ""
	bs, err := generateRandomByteSlice(length)

	if err != nil {
		return "", err
	}

	for _, b := range bs {
		if err != nil {
			return "", err
		}

		c := string(possible[int(b)%possibleLen])
		text += c
	}

	return text, nil
}

func GenerateCodeChallenge(codeVerifier string) string {
	h := sha256.New()
	h.Write([]byte(codeVerifier))
	hashed := h.Sum(nil)

	encoded := b64.URLEncoding.EncodeToString(hashed)
	// needs to remove the '=' character to adhere to Spotify's b64 format
	encoded = strings.ReplaceAll(encoded, "=", "")

	return encoded
}
