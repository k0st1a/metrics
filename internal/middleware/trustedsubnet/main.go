// Package trustedsubnet для проверки подсети, с которой приходят запросы на сервер.
package trustedsubnet

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

type Container interface {
	Contains(string) bool
}

func New(c Container) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ip := r.Header.Get("X-Real-IP")

			if !c.Contains(ip) {
				log.Printf("Request from ip %v is not strusted", ip)
				http.Error(rw, "request ip is not trusted", http.StatusForbidden)
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}
