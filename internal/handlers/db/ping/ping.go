package db

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Pinger interface {
	Ping(context.Context) error
}

type handler struct {
	p Pinger
}

func NewHandler(p Pinger) *handler {
	return &handler{
		p: p,
	}
}

func BuildRouter(r *chi.Mux, h *handler) {
	r.Get("/ping", h.GetPingHandler)
}

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
