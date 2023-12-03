package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/k0st1a/metrics/internal/middleware/logging"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(logging.Logging)

	return r
}
