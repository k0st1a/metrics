// Package db for HTTP ping handler which check connection to DB.
package db

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// Pinger - интерфес проверки доступности БД.
type Pinger interface {
	Ping(context.Context) error
}

type handler struct {
	p Pinger
}

// NewHandler - создание обработчика для проверки доступности БД.
func NewHandler(p Pinger) *handler {
	return &handler{
		p: p,
	}
}

// BuildRouter - формирование маршрута для обработчика.
func BuildRouter(r *chi.Mux, h *handler) {
	r.Get("/ping", h.GetPingHandler)
}

// GetPingHandler - обработчик для проверки доступности БД.
func (h *handler) GetPingHandler(rw http.ResponseWriter, r *http.Request) {
	log.Printf("Get Ping")

	ctx := r.Context()

	err := h.p.Ping(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Ping error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("Get Ping success")

	rw.WriteHeader(http.StatusOK)
}
