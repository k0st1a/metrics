// Package hash is for signature in sha256 format and signature verification.
package hash

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
)

type hash struct {
	key []byte
}

// New - создания сущности подпись.
func New(key string) *hash {
	return &hash{
		key: []byte(key),
	}
}

// Sign - подпись в формате sha256.
func (h *hash) Sign(data []byte) []byte {
	// подписываем алгоритмом HMAC, используя SHA-256
	h1 := hmac.New(sha256.New, h.key)
	h1.Write(data)
	return h1.Sum(nil)
}

// Check - проверка подписи в формате sha256.
func (h *hash) Check(data []byte, sign []byte) bool {
	s := h.Sign(data)
	return bytes.Equal(s, sign)
}
