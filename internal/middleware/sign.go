package middleware

import (
	"bytes"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/k0st1a/metrics/internal/utils"
	"github.com/rs/zerolog/log"
)

type signRoundTrip struct {
	next http.RoundTripper
	hash utils.Signer
}

func NewSign(next http.RoundTripper, sign utils.Signer) *signRoundTrip {
	return &signRoundTrip{
		next: next,
		hash: sign,
	}
}

// Как сделать как в internal/middleware/check_signature.go ?
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
	return s.next.RoundTrip(r)
}
