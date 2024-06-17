// Package middleware логирования запросов на стороне сервера.
package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type responseData struct {
	statusCode  int
	contentSize int
}

type logging struct {
	rw http.ResponseWriter
	rd *responseData
}

func NewLoggingResponse(r http.ResponseWriter) *logging {
	return &logging{
		rw: r,
		rd: &responseData{},
	}
}

func (lr logging) Header() http.Header {
	return lr.rw.Header()
}

func (lr logging) Write(data []byte) (int, error) {
	size, err := lr.rw.Write(data)
	lr.rd.contentSize = size
	if err != nil {
		return size, fmt.Errorf("lr.rw.Write error:%w", err)
	}
	return size, nil
}

func (lr logging) WriteHeader(statusCode int) {
	lr.rw.WriteHeader(statusCode)
	lr.rd.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	logFn := func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lr := NewLoggingResponse(rw)

		next.ServeHTTP(lr, r)

		log.Info().
			Str("uri", r.RequestURI).
			Str("method", r.Method).
			Dur("duration", time.Since(start)).
			Int("status", lr.rd.statusCode).
			Int("size", lr.rd.contentSize).
			Msg("Logging info")
	}
	return http.HandlerFunc(logFn)
}
