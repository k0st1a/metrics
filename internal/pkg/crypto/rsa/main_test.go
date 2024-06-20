package rsa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testBigData = `TEST BIG DATA
MIIBtzCCASwGByqGSM44BAEwggEfAoGBAJ4vZpJ9H6iJR/UU1gJbHTR6in8oa4vX
1Vdvj/V53Q1U2lS0VdkAZyZQiWfO9QTO5oM0Y4S7DtTX3UIiuSuKVWMD55piWuTg
Demf4ZsVAdxcQ6RKCYSwiO0o3O+7RwX2aEzb/KaMqphoHtwRPWhxp5Mbz9kzDD9T
+xQAzsfsuhGVAhUA1kA8zoR9/NuIDs07OdP76UX3UnkCgYEAmB2kVCBqooudn/zU
0dFeXY8RD2OoobKbvdnFeyl8qG3BskLp+1qzHEVT9zI8+6DmJnSxcxyjuT+/ZO1J
nUSX9GNPfWwA4khntera6cLe8qm3fJiWRdsen5XZFFYqvj8A6e5x6qdVCehLGc1Z
Ln0ewTtLDYYpTM/QqFYI7XxKDaEDgYQAAoGAfaNoVmXDVAYAaadSpWgtBuHQNimb
DqOqVQUyEaITd22YMktkccXgwK2XDr4MJT1aBhnIpgpqQ2u+N+EF3JdyxTCtFdKb
PgOIF8OiWe2FjlgoMncOz7SLetQ3f6Y4avpjingyyRwLbDLnEpzSw1fp/v6i0KWL
MIIBtzCCASwGByqGSM44BAEwggEfAoGBAJ4vZpJ9H6iJR/UU1gJbHTR6in8oa4vX
1Vdvj/V53Q1U2lS0VdkAZyZQiWfO9QTO5oM0Y4S7DtTX3UIiuSuKVWMD55piWuTg
Demf4ZsVAdxcQ6RKCYSwiO0o3O+7RwX2aEzb/KaMqphoHtwRPWhxp5Mbz9kzDD9T
+xQAzsfsuhGVAhUA1kA8zoR9/NuIDs07OdP76UX3UnkCgYEAmB2kVCBqooudn/zU
0dFeXY8RD2OoobKbvdnFeyl8qG3BskLp+1qzHEVT9zI8+6DmJnSxcxyjuT+/ZO1J
nUSX9GNPfWwA4khntera6cLe8qm3fJiWRdsen5XZFFYqvj8A6e5x6qdVCehLGc1Z
Ln0ewTtLDYYpTM/QqFYI7XxKDaEDgYQAAoGAfaNoVmXDVAYAaadSpWgtBuHQNimb
DqOqVQUyEaITd22YMktkccXgwK2XDr4MJT1aBhnIpgpqQ2u+N+EF3JdyxTCtFdKb
PgOIF8OiWe2FjlgoMncOz7SLetQ3f6Y4avpjingyyRwLbDLnEpzSw1fp/v6i0KWL
1Vdvj/V53Q1U2lS0VdkAZyZQiWfO9QTO5oM0Y4S7DtTX3UIiuSuKVWMD55piWuTg
Demf4ZsVAdxcQ6RKCYSwiO0o3O+7RwX2aEzb/KaMqphoHtwRPWhxp5Mbz9kzDD9T
+xQAzsfsuhGVAhUA1kA8zoR9/NuIDs07OdP76UX3UnkCgYEAmB2kVCBqooudn/zU
0dFeXY8RD2OoobKbvdnFeyl8qG3BskLp+1qzHEVT9zI8+6DmJnSxcxyjuT+/ZO1J
Ln0ewTtLDYYpTM/QqFYI7XxKDaEDgYQAAoGAfaNoVmXDVAYAaadSpWgtBuHQNimb
DqOqVQUyEaITd22YMktkccXgwK2XDr4MJT1aBhnIpgpqQ2u+N+EF3JdyxTCtFdKb
PgOIF8OiWe2FjlgoMncOz7SLetQ3f6Y4avpjingyyRwLbDLnEpzSw1fp/v6i0KWL
DqOqVQUyEaITd22YMktkccXgwK2XDr4MJT1aBhnIpgpqQ2u+N+EF3JdyxTCtFdKb
PgOIF8OiWe2FjlgoMncOz7SLetQ3f6Y4avpjingyyRwLbDLnEpzSw1fp/v6i0KWL
1Vdvj/V53Q1U2lS0VdkAZyZQiWfO9QTO5oM0Y4S7DtTX3UIiuSuKVWMD55piWuTg
Demf4ZsVAdxcQ6RKCYSwiO0o3O+7RwX2aEzb/KaMqphoHtwRPWhxp5Mbz9kzDD9T
+xQAzsfsuhGVAhUA1kA8zoR9/NuIDs07OdP76UX3UnkCgYEAmB2kVCBqooudn/zU
0dFeXY8RD2OoobKbvdnFeyl8qG3BskLp+1qzHEVT9zI8+6DmJnSxcxyjuT+/ZO1J
nUSX9GNPfWwA4khntera6cLe8qm3fJiWRdsen5XZFFYqvj8A6e5x6qdVCehLGc1Z
Ln0ewTtLDYYpTM/QqFYI7XxKDaEDgYQAAoGAfaNoVmXDVAYAaadSpWgtBuHQNimb
DqOqVQUyEaITd22YMktkccXgwK2XDr4MJT1aBhnIpgpqQ2u+N+EF3JdyxTCtFdKb
PgOIF8OiWe2FjlgoMncOz7SLetQ3f6Y4avpjingyyRwLbDLnEpzSw1fp/v6i0KWL
1Vdvj/V53Q1U2lS0VdkAZyZQiWfO9QTO5oM0Y4S7DtTX3UIiuSuKVWMD55piWuTg
Demf4ZsVAdxcQ6RKCYSwiO0o3O+7RwX2aEzb/KaMqphoHtwRPWhxp5Mbz9kzDD9T
+xQAzsfsuhGVAhUA1kA8zoR9/NuIDs07OdP76UX3UnkCgYEAmB2kVCBqooudn/zU
0dFeXY8RD2OoobKbvdnFeyl8qG3BskLp+1qzHEVT9zI8+6DmJnSxcxyjuT+/ZO1J
nUSX9GNPfWwA4khntera6cLe8qm3fJiWRdsen5XZFFYqvj8A6e5x6qdVCehLGc1Z
Ln0ewTtLDYYpTM/QqFYI7XxKDaEDgYQAAoGAfaNoVmXDVAYAaadSpWgtBuHQNimb
DqOqVQUyEaITd22YMktkccXgwK2XDr4MJT1aBhnIpgpqQ2u+N+EF3JdyxTCtFdKb
PgOIF8OiWe2FjlgoMncOz7SLetQ3f6Y4avpjingyyRwLbDLnEpzSw1fp/v6i0KWL
Tr3hSviZVS0fgEc=`

func TestCrypto(t *testing.T) {
	// to generate file use:
	// openssl rsa -in private.pem -outform PEM -pubout -out public.pem
	pbl, err := NewPublicFromFile("./public.pem")
	assert.NoError(t, err)

	// to generate file use:
	// openssl genrsa -out private.pem 4096
	prv, err := NewPrivateFromFile("./private.pem")
	assert.NoError(t, err)

	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "check encrypt and decrypt",
			data: []byte("Hello Gophers"),
		},
		{
			name: "check encrypt and decrypt big data",
			data: []byte(testBigData),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encData, err := pbl.Encrypt(test.data)
			assert.NoError(t, err)

			decData, err := prv.Decrypt(encData)
			assert.NoError(t, err)

			assert.Equal(t, test.data, decData)
		})
	}
}

func TestPublicNoPEMData(t *testing.T) {
	pbl, err := NewPublic("no PEM data")
	assert.Nil(t, pbl)
	assert.EqualError(t, err, "no PEM data is found")
}

func TestPublicKeyParseError(t *testing.T) {
	badPublicKey := "-----BEGIN PUBLIC KEY-----\n" +
		"MIIJKAIBAAKCAgEAo5i9CqGMbZw1dS8Jx7Ne9/9SYZBcOxe/39Dbwz+oL1jMpoPJ\n" +
		"wI4YOkfNC7c2+TozI00C9a2KbHnk2L3sSBNUXL875YHqO3tnz2Uz5Vvvewti6Lcl\n" +
		"-----END PUBLIC KEY-----"
	prv, err := NewPublic(badPublicKey)
	assert.Nil(t, prv)
	assert.ErrorContains(t, err, "x509 pkix public key parse error")
}

func TestPublicKeyBadTypeAssertion(t *testing.T) {
	var testDSAPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBtzCCASwGByqGSM44BAEwggEfAoGBAJ4vZpJ9H6iJR/UU1gJbHTR6in8oa4vX
1Vdvj/V53Q1U2lS0VdkAZyZQiWfO9QTO5oM0Y4S7DtTX3UIiuSuKVWMD55piWuTg
Demf4ZsVAdxcQ6RKCYSwiO0o3O+7RwX2aEzb/KaMqphoHtwRPWhxp5Mbz9kzDD9T
+xQAzsfsuhGVAhUA1kA8zoR9/NuIDs07OdP76UX3UnkCgYEAmB2kVCBqooudn/zU
0dFeXY8RD2OoobKbvdnFeyl8qG3BskLp+1qzHEVT9zI8+6DmJnSxcxyjuT+/ZO1J
nUSX9GNPfWwA4khntera6cLe8qm3fJiWRdsen5XZFFYqvj8A6e5x6qdVCehLGc1Z
Ln0ewTtLDYYpTM/QqFYI7XxKDaEDgYQAAoGAfaNoVmXDVAYAaadSpWgtBuHQNimb
DqOqVQUyEaITd22YMktkccXgwK2XDr4MJT1aBhnIpgpqQ2u+N+EF3JdyxTCtFdKb
PgOIF8OiWe2FjlgoMncOz7SLetQ3f6Y4avpjingyyRwLbDLnEpzSw1fp/v6i0KWL
Tr3hSviZVS0fgEc=
-----END PUBLIC KEY-----`

	prv, err := NewPublic(testDSAPublicKey)
	assert.Nil(t, prv)
	assert.ErrorContains(t, err, "bad type assertion to *rsa.PublicKey")
}

func TestPubliceReadFromFile(t *testing.T) {
	prv, err := NewPublicFromFile("bad path to file")
	assert.Nil(t, prv)
	assert.ErrorContains(t, err, "public key file read error")
}

func TestPublicEncrypt(t *testing.T) {
	pbl, err := NewPublicFromFile("./public.pem")
	assert.NoError(t, err)

	pbl.key.E = -1
	data, err := pbl.Encrypt([]byte("some bad encrypted data"))
	assert.Nil(t, data)
	assert.ErrorContains(t, err, "rsa OAEP encrypt error")
}

func TestPrivateNoPEMData(t *testing.T) {
	prv, err := NewPrivate("no PEM data")
	assert.Nil(t, prv)
	assert.EqualError(t, err, "no PEM data is found")
}

func TestPrivateKeyParseError(t *testing.T) {
	badPrivateKey := "-----BEGIN RSA PRIVATE KEY-----\n" +
		"MIIJKAIBAAKCAgEAo5i9CqGMbZw1dS8Jx7Ne9/9SYZBcOxe/39Dbwz+oL1jMpoPJ\n" +
		"wI4YOkfNC7c2+TozI00C9a2KbHnk2L3sSBNUXL875YHqO3tnz2Uz5Vvvewti6Lcl\n" +
		"-----END RSA PRIVATE KEY-----"
	prv, err := NewPrivate(badPrivateKey)
	assert.Nil(t, prv)
	assert.ErrorContains(t, err, "x509 pkcs1 private key parse error")
}

func TestPrivateReadFromFile(t *testing.T) {
	prv, err := NewPrivateFromFile("bad path to file")
	assert.Nil(t, prv)
	assert.ErrorContains(t, err, "private key file read error")
}

func TestPrivateDecrypt(t *testing.T) {
	prv, err := NewPrivateFromFile("./private.pem")
	assert.NoError(t, err)

	data, err := prv.Decrypt([]byte("some bad encrypted data"))
	assert.Nil(t, data)
	assert.ErrorContains(t, err, "rsa decrypt OAEP error")
}
