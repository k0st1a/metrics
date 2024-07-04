// Package handler for create HTTP router.
package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// NewRouter - создание нового маршрутизатора.
func NewRouter(middlewares []func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares...)
	return r
}
