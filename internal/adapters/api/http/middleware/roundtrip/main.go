// Package roundtrip for chaining middlewares for http client.
package roundtrip

import (
	"net/http"
)

type HandlerFunc func(*http.Request) (*http.Response, error)

func (f HandlerFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

type Middleware func(http.RoundTripper) http.RoundTripper

func New(h http.RoundTripper, middlewares ...Middleware) http.RoundTripper {
	for _, m := range middlewares {
		h = m(h)
	}

	return h
}
