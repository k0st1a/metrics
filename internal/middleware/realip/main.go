// Package realip для выставления заголовка X-Real-IP.
package realip

import (
	"net/http"

	"github.com/k0st1a/metrics/internal/middleware/roundtrip"
)

func New(ip string) func(http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return roundtrip.HandlerFunc(func(r *http.Request) (*http.Response, error) {
			r.Header.Set("X-Real-IP", ip)
			return next.RoundTrip(r)
		})
	}
}
