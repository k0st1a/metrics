// Package encrypt для кодирования данных, отправляемых со стороны агента.
package encrypt

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/k0st1a/metrics/internal/middleware/roundtrip"
)

type Encrypter interface {
	Encrypt([]byte) ([]byte, error)
}

func New(enc Encrypter) func(http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return roundtrip.HandlerFunc(func(r *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return nil, fmt.Errorf("body read error while encrypt:%w", err)
			}

			err = r.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("body close error while encrypt:%w", err)
			}

			encBody, err := enc.Encrypt(body)
			if err != nil {
				return nil, fmt.Errorf("encrypt body error:%w", err)
			}

			r.Body = io.NopCloser(bytes.NewBuffer(encBody))

			r.ContentLength = int64(len(encBody))

			//nolint:wrapcheck //no need here
			return next.RoundTrip(r)
		})
	}
}
