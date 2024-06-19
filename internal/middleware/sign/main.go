// Package sign is middleware for sign sending data from HTTP client.
package sign

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"github.com/k0st1a/metrics/internal/middleware/roundtrip"
)

type Signer interface {
	Sign([]byte) []byte
}

func New(s Signer) func(http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return roundtrip.HandlerFunc(func(r *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return nil, fmt.Errorf("body read error while sign")
			}

			err = r.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("body close error while sign")
			}

			r.Body = io.NopCloser(bytes.NewBuffer(body))

			signBody := s.Sign(body)
			hex := hex.EncodeToString(signBody)
			r.Header.Set("HashSHA256", hex)

			//nolint:wrapcheck //no need here
			return next.RoundTrip(r)
		})
	}
}
