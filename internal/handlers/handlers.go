package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/k0st1a/metrics/internal/middleware"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logging)
	r.Use(middleware.Compress)

	return r
}
