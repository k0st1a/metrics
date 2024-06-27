// Package trustedsubnet для проверки подсети, с которой приходят запросы на сервер.
package trustedsubnet

import (
	"net/http"
)

type Container interface {
	Contains(string) bool
}

func New(c Container) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			ip := r.Header.Get("X-Real-IP")

			if !c.Contains(ip) {
				http.Error(rw, "request ip is not trusted", http.StatusForbidden)
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}
