package middleware

import (
	"bytes"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/k0st1a/metrics/internal/pkg/hash"
	"github.com/rs/zerolog/log"
)

type signRoundTrip struct {
	next http.RoundTripper
	hash hash.Signer
}

// NewSign - `middleware` для подписи отправляемых с клиента данных.
func NewSign(next http.RoundTripper, sign hash.Signer) *signRoundTrip {
	return &signRoundTrip{
		next: next,
		hash: sign,
	}
}

// RoundTrip - вызывается на клиенте перед отправкой данных на сервер.
func (s *signRoundTrip) RoundTrip(r *http.Request) (*http.Response, error) {
	if s.hash.Is() {
		b, err := io.ReadAll(r.Body)
		cerr := r.Body.Close()
		if cerr != nil {
			log.Error().Err(err).Msg("body close error while sign")
		}
		if err != nil {
			log.Error().Err(err).Msg("body read error while sign")
		}
		r.Body = io.NopCloser(bytes.NewBuffer(b))

		s := s.hash.Sign(b)
		hs := hex.EncodeToString(s)
		r.Header.Set("HashSHA256", hs)
	}
	//nolint:wrapcheck //no need here
	return s.next.RoundTrip(r)
}
