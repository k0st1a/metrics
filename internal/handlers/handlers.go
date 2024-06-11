package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/k0st1a/metrics/internal/middleware"
	"github.com/k0st1a/metrics/internal/pkg/hash"
)

// NewRouter - создание нового маршрутизатора.
func NewRouter(h hash.Checker) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logging)
	r.Use(middleware.Compress)
	r.Use(middleware.CheckSignature(h))

	return r
}
