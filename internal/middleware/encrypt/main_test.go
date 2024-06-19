package encrypt

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/k0st1a/metrics/internal/middleware/decrypt"
	"github.com/k0st1a/metrics/internal/middleware/roundtrip"
	"github.com/k0st1a/metrics/internal/pkg/crypto/rsa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type errReader int

func (errReader) Read(p []byte) (int, error) {
	return 0, errors.New("test read error")
}

type errCloser int

func (errCloser) Read(p []byte) (int, error) {
	return 0, io.EOF
}
func (errCloser) Close() error {
	return errors.New("test close error")
}

type errEncrypter int

func (errEncrypter) Encrypt(_ []byte) ([]byte, error) {
	return nil, errors.New("test encrypt error")
}

var responseRoundTripper http.RoundTripper = testRoundTripper(0)

type testRoundTripper int

func (testRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return &http.Response{
		Body:   r.Body,
		Header: r.Header,
	}, nil
}

func TestEncryptError(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		doError string
		body    io.Reader
		encrypt Encrypter
	}{
		{
			name:    "check body read error",
			file:    "./public.pem",
			body:    errReader(0),
			doError: "body read error while encrypt",
		},
		{
			name:    "check body close error",
			file:    "./public.pem",
			body:    errCloser(0),
			doError: "body close error while encrypt",
		},
		{
			name:    "check data encrypt error",
			body:    bytes.NewBuffer([]byte("good body")),
			encrypt: errEncrypter(0),
			doError: "encrypt body error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testServer := httptest.NewServer(nil)
			defer testServer.Close()

			req, err := http.NewRequest(http.MethodPost, testServer.URL, test.body)
			assert.NoError(t, err)

			var h Encrypter

			if test.file != "" {
				h, err = rsa.NewPublicFromFile(test.file)
				assert.NoError(t, err)
			} else if test.encrypt != nil {
				h = test.encrypt
			}

			rt := roundtrip.New(responseRoundTripper, New(h))

			c := &http.Client{
				Transport: rt,
			}

			resp, err := c.Do(req)
			assert.Nil(t, resp)
			assert.ErrorContains(t, err, test.doError)

			if resp != nil {
				err = resp.Body.Close()
				assert.NoError(t, err)
			}
		})
	}
}

func TestEncryptAndDecrypt(t *testing.T) {
	tests := []struct {
		name                  string
		publicKeyFile         string
		privateKeyFile        string
		body                  string
		expectedContentLength int
	}{
		{
			name:                  "check encrypt and decrypt",
			publicKeyFile:         "./public.pem",
			privateKeyFile:        "../decrypt/private.pem",
			body:                  "кодируемые данные",
			expectedContentLength: 33,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pbl, err := rsa.NewPublicFromFile(test.publicKeyFile)
			require.NoError(t, err)

			prv, err := rsa.NewPrivateFromFile(test.privateKeyFile)
			require.NoError(t, err)

			r := chi.NewRouter()
			r.Use(decrypt.New(prv))

			r.Post("/check", func(rw http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				err = r.Body.Close()
				require.NoError(t, err)

				cl, err := strconv.Atoi(r.Header.Get("Content-Length"))
				require.NoError(t, err)

				assert.Equal(t, test.expectedContentLength, cl)

				assert.Equal(t, []byte(test.body), body)

				rw.WriteHeader(http.StatusOK)
			})

			testServer := httptest.NewServer(r)
			defer testServer.Close()

			req, err := http.NewRequest(http.MethodPost, testServer.URL+"/check", bytes.NewBuffer([]byte(test.body)))
			assert.NoError(t, err)

			rt := roundtrip.New(http.DefaultTransport, New(pbl))
			c := &http.Client{
				Transport: rt,
			}

			resp, err := c.Do(req)
			assert.NoError(t, err)

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			err = resp.Body.Close()
			assert.NoError(t, err)
		})
	}
}
