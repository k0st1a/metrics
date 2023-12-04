package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

type compress struct {
	rw http.ResponseWriter
	w  io.WriteCloser
}

func newCompress(rw http.ResponseWriter) *compress {
	return &compress{
		rw: rw,
		w:  gzip.NewWriter(rw),
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
	if !isNeedCompress(ct) {
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
		ct := r.Header.Get("Content-Type")

		ac := r.Header.Get("Accept-Encoding")
		if isNeedCompress(ct) && ac == "gzip" {
			c := newCompress(rw)
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
