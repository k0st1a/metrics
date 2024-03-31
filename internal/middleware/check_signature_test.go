package middleware

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/k0st1a/metrics/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestCheckSignature(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		sign     string
		body     io.Reader
		want     int
		wantBody string
	}{
		{
			name:     "Ошибка декодирования подписи",
			key:      "some key",
			sign:     "bad sign",
			want:     400,
			wantBody: "hash decode error\n",
		},
		{
			name:     "Ошибка чтения тела",
			key:      "some key",
			sign:     "2a2629ba328d5376b44f88536047a12500d33bc43045a7407c29a88312bc2a48",
			body:     errReader(0),
			want:     400,
			wantBody: "body read error\n",
		},
		{
			name:     "Подпись в запросе не верная",
			key:      "some key",
			sign:     "2a2629ba328d5376b44f88536047a12500d33bc43045a7407c29a88312bc2a48",
			body:     bytes.NewBuffer([]byte{}),
			want:     400,
			wantBody: "wrong signature\n",
		},
		{
			name:     "Подпись в запросе верная",
			key:      "some key",
			sign:     "e38c1bd0a6f6196624b914c454929f684c19ffbe0b8d59954bcb0498e17cc165",
			body:     bytes.NewBuffer([]byte("подписываемые данные")),
			want:     200,
			wantBody: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			h := utils.NewHash(test.key)

			r := chi.NewRouter()
			r.Use(CheckSignature(h))
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {})

			req := httptest.NewRequest(http.MethodPost, "/", test.body)
			req.Header.Set("HashSHA256", test.sign)

			r.ServeHTTP(recorder, req)
			res := recorder.Result()

			require.Equal(t, test.want, res.StatusCode)

			b, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			err = res.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, test.wantBody, string(b))
		})
	}
}
