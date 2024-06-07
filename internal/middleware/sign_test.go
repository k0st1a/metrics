package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/k0st1a/metrics/internal/pkg/hash"

	"github.com/stretchr/testify/assert"
)

func TestCheckSign(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		sign     string
		body     io.Reader
		want     int
		wantBody string
	}{
		{
			name:     "Sign decode error",
			key:      "some key",
			sign:     "bad sign",
			body:     bytes.NewBuffer([]byte("подписываемые данные")),
			want:     400,
			wantBody: "hash decode error\n",
		},
		{
			name:     "Sign in request not correct",
			key:      "some key",
			sign:     "2a2629ba328d5376b44f88536047a12500d33bc43045a7407c29a88312bc2a48",
			body:     bytes.NewBuffer([]byte{}),
			want:     400,
			wantBody: "wrong signature\n",
		},
		{
			name:     "Sign is correct",
			key:      "some key",
			sign:     "e38c1bd0a6f6196624b914c454929f684c19ffbe0b8d59954bcb0498e17cc165",
			body:     bytes.NewBuffer([]byte("подписываемые данные")),
			want:     200,
			wantBody: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := hash.New(test.key)

			r := chi.NewRouter()
			r.Use(CheckSignature(h))
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {})

			testServer := httptest.NewServer(r)
			defer testServer.Close()

			req, err := http.NewRequest(http.MethodPost, testServer.URL, test.body)
			assert.NoError(t, err)

			crt := NewChangeRoundTrip(http.DefaultTransport, test.sign)
			sign := NewSign(crt, h)

			c := &http.Client{
				Transport: sign,
			}

			resp, err := c.Do(req)
			assert.NoError(t, err)

			respBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			err = resp.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, test.wantBody, string(respBody))
		})
	}
}

type changeRoundTrip struct {
	next http.RoundTripper
	sign string
}

func NewChangeRoundTrip(next http.RoundTripper, sign string) *changeRoundTrip {
	return &changeRoundTrip{
		next: next,
		sign: sign,
	}
}

func (s *changeRoundTrip) RoundTrip(r *http.Request) (*http.Response, error) {
	b, err := io.ReadAll(r.Body)
	cerr := r.Body.Close()
	if cerr != nil {
		return nil, fmt.Errorf("body close error:%w", cerr)
	}
	if err != nil {
		return nil, fmt.Errorf("body read error:%w", err)
	}
	r.Body = io.NopCloser(bytes.NewBuffer(b))

	r.Header.Set("HashSHA256", s.sign)

	//nolint:wrapcheck //no need here
	return s.next.RoundTrip(r)
}
