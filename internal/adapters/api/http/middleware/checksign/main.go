// Package checksign for check signature HashSHA256 for incoming msg of HTTP server.
package checksign

import (
	"bytes"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

// Checker - интерфейс проверки подписи данных.
type Checker interface {
	Check(data []byte, sign []byte) (equal bool)
}

func New(h Checker) func(next http.Handler) http.Handler {
	// Подсмотрел в https://github.com/go-chi/chi/blob/master/middleware/content_type.go
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			sign := r.Header.Get("HashSHA256")
			if sign != "" {
				ds, err := hex.DecodeString(sign)
				if err != nil {
					log.Error().Err(err).Msg("hash decode error while checksign")
					http.Error(rw, "hash decode error while checksign", http.StatusBadRequest)
					return
				}

				b, err := io.ReadAll(r.Body)
				if err != nil {
					log.Error().Err(err).Msg("body read error while checksign")
					http.Error(rw, "body read error while checksign", http.StatusBadRequest)
					return
				}

				err = r.Body.Close()
				if err != nil {
					log.Error().Err(err).Msg("body close error while checksign")
				}

				if !h.Check(b, ds) {
					log.Error().Err(err).Msg("wrong signature")
					http.Error(rw, "wrong signature", http.StatusBadRequest)
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
