package decrypt

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
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

type errDecrypter int

func (errDecrypter) Decrypt(_ []byte) ([]byte, error) {
	return nil, errors.New("test decrypt error")
}

func TestDecryptError(t *testing.T) {
	tests := []struct {
		name           string
		file           string
		body           io.Reader
		decrypt        Decrypter
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name:           "check body read error",
			file:           "./private.pem",
			body:           errReader(0),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   []byte("body read error while decrypt\n"),
		},
		{
			name:           "check body close error",
			body:           errCloser(0),
			decrypt:        errDecrypter(0),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   []byte("decrypt body error\n"),
		},
		{
			name:           "check data encrypt error",
			body:           bytes.NewBuffer([]byte("good body")),
			decrypt:        errDecrypter(0),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   []byte("decrypt body error\n"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				h   Decrypter
				err error
			)

			if test.file != "" {
				h, err = rsa.NewPrivateFromFile(test.file)
				require.NoError(t, err)
			} else if test.decrypt != nil {
				h = test.decrypt
			}

			r := chi.NewRouter()
			r.Use(New(h))
			w := httptest.NewRecorder()

			r.Post("/check", func(rw http.ResponseWriter, r *http.Request) {})

			testServer := httptest.NewServer(r)
			defer testServer.Close()

			req, err := http.NewRequest(http.MethodPost, testServer.URL+"/check", test.body)
			assert.NoError(t, err)

			r.ServeHTTP(w, req)
			resp := w.Result()

			require.Equal(t, test.expectedStatus, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, test.expectedBody, body)

			err = resp.Body.Close()
			assert.NoError(t, err)
		})
	}
}
