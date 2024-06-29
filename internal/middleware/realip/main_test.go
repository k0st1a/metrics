package realip

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k0st1a/metrics/internal/middleware/roundtrip"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var responseRoundTripper http.RoundTripper = testRoundTripper(0)

type testRoundTripper int

func (testRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return &http.Response{
		Body:   r.Body,
		Header: r.Header,
	}, nil
}

func TestRealIP(t *testing.T) {
	tests := []struct {
		name    string
		xRealIP string
	}{
		{
			name:    "check set header X-Real-IP in request",
			xRealIP: "127.0.0.1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testServer := httptest.NewServer(nil)
			defer testServer.Close()

			req, err := http.NewRequest(http.MethodPost, testServer.URL, nil)
			assert.NoError(t, err)

			rt := roundtrip.New(responseRoundTripper, New(test.xRealIP))
			c := &http.Client{
				Transport: rt,
			}

			resp, err := c.Do(req)
			assert.NoError(t, err)

			err = resp.Body.Close()
			assert.NoError(t, err)

			require.Equal(t, test.xRealIP, resp.Header.Get("X-Real-IP"))
		})
	}
}
