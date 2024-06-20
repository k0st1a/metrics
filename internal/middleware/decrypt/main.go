// Package decrypt для декодирования сообщений, обрабатываемых сервером.
package decrypt

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

type Decrypter interface {
	Decrypt([]byte) ([]byte, error)
}

func New(d Decrypter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				log.Error().Err(err).Msg("body read error while decrypt")
				http.Error(rw, "body read error while decrypt", http.StatusBadRequest)
				return
			}

			err = r.Body.Close()
			if err != nil {
				log.Error().Err(err).Msg("body close error while decrypt")
			}

			dec, err := d.Decrypt(b)
			if err != nil {
				log.Error().Err(err).Msg("decrypt body error")
				http.Error(rw, "decrypt body error", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(dec))

			cl := len(dec)
			r.ContentLength = int64(cl)
			r.Header.Set("Content-Length", strconv.Itoa(cl))

			next.ServeHTTP(rw, r)
		})
	}
}
