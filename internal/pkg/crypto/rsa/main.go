// Package rsa is for data rsa encryption.
package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"hash"
	"os"
)

type private struct {
	key  *rsa.PrivateKey
	hash hash.Hash
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
		key:  key,
		hash: sha256.New(),
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
	var decData []byte

	dataLen := len(data)
	step := p.chunkSize()
	for begin := 0; begin < dataLen; begin += step {
		end := begin + step
		if end > dataLen {
			end = dataLen
		}

		decChunk, err := p.decryptChunk(data[begin:end])
		if err != nil {
			return nil, fmt.Errorf("decrypt chunk error:%w", err)
		}

		decData = append(decData, decChunk...)
	}

	return decData, nil
}

func (p *private) decryptChunk(data []byte) ([]byte, error) {
	b, err := rsa.DecryptOAEP(p.hash, rand.Reader, p.key, data, nil)
	if err != nil {
		return nil, fmt.Errorf("rsa decrypt OAEP error:%w", err)
	}

	return b, nil
}

func (p *private) chunkSize() int {
	return p.key.Size()
}

type public struct {
	key  *rsa.PublicKey
	hash hash.Hash
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
		key:  key,
		hash: sha256.New(),
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
	var encData []byte

	dataLen := len(data)
	step := p.chunkSize()
	for begin := 0; begin < dataLen; begin += step {
		end := begin + step
		if end > dataLen {
			end = dataLen
		}

		encChunk, err := p.encryptChunk(data[begin:end])
		if err != nil {
			return nil, fmt.Errorf("encrypt chunk error:%w", err)
		}

		encData = append(encData, encChunk...)
	}

	return encData, nil
}

func (p *public) encryptChunk(data []byte) ([]byte, error) {
	b, err := rsa.EncryptOAEP(p.hash, rand.Reader, p.key, data, nil)
	if err != nil {
		return nil, fmt.Errorf("rsa OAEP encrypt error:%w", err)
	}

	return b, nil
}

// The message must be no longer than the length of the public modulus minus
// twice the hash length, minus a further 2.
func (p *public) chunkSize() int {
	return p.key.Size() - 2*p.hash.Size() - 2
}
