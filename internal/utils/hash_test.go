package utils

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashSign(t *testing.T) {
	tests := []struct {
		name string
		key  string
		data []byte
		sign string
	}{
		{
			name: "Проверка подписи SHA256",
			key:  "какой-то ключ",
			data: []byte("подписываемые данные"),
			sign: "22fa39b2b38f02f3f922e651d7635b27ac023dd20ae05b9c758cfe1d3044e831",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := NewHash(test.key)

			ds, err := hex.DecodeString(test.sign)
			assert.NoError(t, err)

			assert.Equal(t, ds, h.Sign(test.data))

			assert.True(t, h.CheckSignature(test.data, ds))
		})
	}
}

func TestHashIs(t *testing.T) {
	tests := []struct {
		name string
		key  string
		is   bool
	}{
		{
			name: "Проверка hash с не пустым ключом",
			key:  "какой-то ключ",
			is:   true,
		},
		{
			name: "Проверка hash с пустым ключом",
			key:  "",
			is:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := NewHash(test.key)

			assert.Equal(t, test.is, h.Is())
		})
	}
}
