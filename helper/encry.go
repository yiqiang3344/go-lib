package helper

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func ComputeHmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sha := h.Sum(nil)
	return base64.StdEncoding.EncodeToString([]byte(sha))
}
