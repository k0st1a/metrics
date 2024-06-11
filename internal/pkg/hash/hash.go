// Пакет подписи в формате sha256 и проверки подписи.
package hash

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
)

// Signer - интерфейс подписи данных.
type Signer interface {
	Sign(data []byte) (sign []byte)
	Is() bool
}

// Checker - интерфейс проверки подписи данных.
type Checker interface {
	Check(data []byte, sign []byte) (equal bool)
	Is() bool
}

type hash struct {
	key []byte
}

// New - создания сущности подпись.
func New(key string) *hash {
	return &hash{
		key: []byte(key),
	}
}

// Is - проверка есть ли подпись.
func (h *hash) Is() bool {
	return !bytes.Equal(h.key, []byte{})
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
