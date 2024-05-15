package hash

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
)

type Signer interface {
	Sign(data []byte) (sign []byte)
	Is() bool
}

type Checker interface {
	Check(data []byte, sign []byte) (equal bool)
	Is() bool
}

type hash struct {
	key []byte
}

func New(key string) *hash {
	return &hash{
		key: []byte(key),
	}
}

func (h *hash) Is() bool {
	return !bytes.Equal(h.key, []byte{})
}

func (h *hash) Sign(data []byte) []byte {
	// подписываем алгоритмом HMAC, используя SHA-256
	h1 := hmac.New(sha256.New, h.key)
	h1.Write(data)
	return h1.Sum(nil)
}

func (h *hash) Check(data []byte, sign []byte) bool {
	s := h.Sign(data)
	return bytes.Equal(s, sign)
}
