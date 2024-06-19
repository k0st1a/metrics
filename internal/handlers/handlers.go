// Package handlers for create HTTP router.
package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// NewRouter - создание нового маршрутизатора.
func NewRouter(middlewares []func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares...)
	return r
}
