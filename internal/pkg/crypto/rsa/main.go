// Package rsa is for data rsa encryption.
package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type private struct {
	key *rsa.PrivateKey
}

// NewPrivate create rsa private side by data of private key.
func NewPrivate(data string) (*private, error) {
	block, _ := pem.Decode([]byte(data))
	if block == nil {
		return nil, fmt.Errorf("no PEM data is found")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("x509 pkcs1 private key parse error:%w", err)
	}

	return &private{
		key: key,
	}, nil
}

// NewPrivateFromFile create rsa private side by file of private key.
func NewPrivateFromFile(path string) (*private, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("private key file read error:%w", err)
	}

	return NewPrivate(string(data))
}

// Decrypt data.
func (p *private) Decrypt(data []byte) ([]byte, error) {
	b, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, p.key, data, nil)
	if err != nil {
		return nil, fmt.Errorf("rsa decrypt OAEP error:%w", err)
	}

	return b, nil
}

type public struct {
	key *rsa.PublicKey
}

// NewPublic create rsa public side by data of public key.
func NewPublic(data string) (*public, error) {
	block, _ := pem.Decode([]byte(data))
	if block == nil {
		return nil, fmt.Errorf("no PEM data is found")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("x509 pkix public key parse error:%w", err)
	}

	key, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("bad type assertion to *rsa.PublicKey")
	}

	return &public{
		key: key,
	}, nil
}

// NewPublicFromFile create rsa public side by file of public key.
func NewPublicFromFile(path string) (*public, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("public key file read error:%w", err)
	}

	return NewPublic(string(data))
}

// Encryp data.
func (p *public) Encrypt(data []byte) ([]byte, error) {
	b, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, p.key, data, nil)
	if err != nil {
		return nil, fmt.Errorf("rsa OAEP encrypt error:%w", err)
	}

	return b, nil
}
