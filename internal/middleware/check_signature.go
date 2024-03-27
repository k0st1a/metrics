package middleware

import (
	"bytes"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/k0st1a/metrics/internal/utils"
	"github.com/rs/zerolog/log"
)

func CheckSignature(h utils.SignChecker) func(next http.Handler) http.Handler {
	// Подсмотрел в https://github.com/go-chi/chi/blob/master/middleware/content_type.go
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			sign := r.Header.Get("HashSHA256")
			if sign != "" {
				ds, err := hex.DecodeString(sign)
				if err != nil {
					log.Error().Err(err).Msg("hash decode error")
					http.Error(rw, "hash decode error", http.StatusInternalServerError)
					return
				}

				b, err := io.ReadAll(r.Body)
				cerr := r.Body.Close()
				if cerr != nil {
					log.Error().Err(err).Msg("body close error")
				}
				if err != nil {
					log.Error().Err(err).Msg("request body read error while check signature")
					http.Error(rw, "body read error", http.StatusInternalServerError)
					return
				}

				if h.Is() && !h.CheckSignature(b, ds) {
					log.Error().Err(err).Msg("unsuccess check signature")
					http.Error(rw, "unsuccess check signature", http.StatusBadRequest)
					return
				}

				// Восстанавливаем тело запроса нашел на stackoverflow:
				// https://stackoverflow.com/questions/46948050/how-to-read-request-body-twice-in-golang-middleware
				r.Body = io.NopCloser(bytes.NewBuffer(b))
			}
			next.ServeHTTP(rw, r)
		})
	}
}
