// Middleware сжатия данных для Content-Type application/json и text/html на стороне сервера.
package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type compress struct {
	rw     http.ResponseWriter
	w      io.WriteCloser
	uri    string
	method string
}

func newCompress(rw http.ResponseWriter, r *http.Request) *compress {
	return &compress{
		rw:     rw,
		w:      gzip.NewWriter(rw),
		uri:    r.RequestURI,
		method: r.Method,
	}
}

func (c *compress) close() {
	_ = c.w.Close()
}

func (c compress) Header() http.Header {
	return c.rw.Header()
}

func (c compress) Write(data []byte) (int, error) {
	ct := c.Header().Get("Content-Type")
	in := isNeedCompress(ct)
	log.Debug().Msgf("Is need compress for Content-Type:%v?(%v) method:%v, uri:%v", ct, in, c.method, c.uri)
	if !in {
		n, err := c.rw.Write(data)
		if err != nil {
			return n, fmt.Errorf("c.rw.Write error:%w", err)
		}
		return n, nil
	}

	c.Header().Set("Content-Encoding", "gzip")

	n, err := c.w.Write(data)
	if err != nil {
		return n, fmt.Errorf("c.w.Write error:%w", err)
	}
	return n, nil
}

func (c compress) WriteHeader(statusCode int) {
	c.rw.WriteHeader(statusCode)
}

func Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept-Encoding") == "gzip" {
			c := newCompress(rw, r)
			next.ServeHTTP(c, r)
			defer func() {
				c.close()
			}()
		} else {
			next.ServeHTTP(rw, r)
		}
	})
}

func isNeedCompress(ct string) bool {
	switch ct {
	case "application/json":
		return true
	case "text/html":
		return true
	default:
		return false
	}
}
