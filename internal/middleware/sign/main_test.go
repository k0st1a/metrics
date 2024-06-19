package sign

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k0st1a/metrics/internal/middleware/roundtrip"
	"github.com/k0st1a/metrics/internal/pkg/hash"

	"github.com/stretchr/testify/assert"
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

var responseRoundTripper http.RoundTripper = testRoundTripper(0)

type testRoundTripper int

func (testRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return &http.Response{
		Body:   r.Body,
		Header: r.Header,
	}, nil
}

func TestSignError(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		doError string
		body    io.Reader
	}{
		{
			name:    "check body read error",
			key:     "some key",
			body:    errReader(0),
			doError: "body read error while sign",
		},
		{
			name:    "check body close error",
			key:     "some key",
			body:    errCloser(0),
			doError: "body close error while sign",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testServer := httptest.NewServer(nil)
			defer testServer.Close()

			req, err := http.NewRequest(http.MethodPost, testServer.URL, test.body)
			assert.NoError(t, err)

			h := hash.New(test.key)
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

func TestSign(t *testing.T) {
	tests := []struct {
		name string
		key  string
		body string
		sign string
	}{
		{
			name: "check set HashSHA256",
			key:  "some key",
			body: "подписываемые данные",
			sign: "e38c1bd0a6f6196624b914c454929f684c19ffbe0b8d59954bcb0498e17cc165",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testServer := httptest.NewServer(nil)
			defer testServer.Close()

			req, err := http.NewRequest(http.MethodPost, testServer.URL, bytes.NewBuffer([]byte(test.body)))
			assert.NoError(t, err)

			h := hash.New(test.key)
			rt := roundtrip.New(responseRoundTripper, New(h))
			c := &http.Client{
				Transport: rt,
			}

			resp, err := c.Do(req)
			assert.NoError(t, err)

			assert.Equal(t, test.sign, resp.Header.Get("HashSHA256"))

			respBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			err = resp.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, test.body, string(respBody))
		})
	}
}
