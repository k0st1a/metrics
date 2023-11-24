package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/k0st1a/metrics/internal/storage/counter"
	"github.com/k0st1a/metrics/internal/storage/gauge"
	"github.com/rs/zerolog/log"
)

type storageService interface {
	GetGauge(string) (float64, bool)
	StoreGauge(string, float64)
	GetCounter(string) (int64, bool)
	StoreCounter(string, int64)
}

type handler struct {
	storage storageService
}

func NewHandler(s storageService) *handler {
	return &handler{
		storage: s,
	}
}

func BuildRouter(h *handler) chi.Router {
	log.Debug().Msg("Make router")
	r := chi.NewRouter()

	log.Debug().Msg("POST handlers adding")
	r.Route("/update/", func(r chi.Router) {
		r.Post("/counter/{name}/{value}", h.PostCounterHandler)
		r.Post("/counter/", h.NotFoundHandler)
		r.Post("/gauge/{name}/{value}", h.PostGaugeHandler)
		r.Post("/gauge/", h.NotFoundHandler)
	})
	log.Debug().Msg("POST handlers added")

	log.Debug().Msg("GET handlers adding")
	r.Route("/value/", func(r chi.Router) {
		r.Get("/counter/{name}", h.GetCounterHandler)
		r.Get("/gauge/{name}", h.GetGaugeHandler)
	})
	log.Debug().Msg("GET handlers added")

	log.Debug().Msg("Custom NotFound handler adding")
	r.NotFound(h.BadRequestHandler)
	log.Debug().Msg("Custom NotFound handler added")

	return r
}

func (h *handler) BadRequestHandler(rw http.ResponseWriter, r *http.Request) {
	log.Debug().
		Str("RequestURI", r.RequestURI).
		Msg("")

	rw.WriteHeader(http.StatusBadRequest)
}

func (h *handler) NotFoundHandler(rw http.ResponseWriter, r *http.Request) {
	log.Debug().
		Str("RequestURI", r.RequestURI).
		Msg("")

	http.Error(rw, "metric value is empty", http.StatusNotFound)
}

func (h *handler) PostCounterHandler(rw http.ResponseWriter, r *http.Request) {
	name := strings.ToLower(chi.URLParam(r, "name"))
	log.Debug().
		Str("RequestURI", r.RequestURI).
		Str("name", name).
		Msg("")

	if name == "" {
		http.Error(rw, "metric name is empty", http.StatusNotFound)
		return
	}

	value := strings.ToLower(chi.URLParam(r, "value"))
	if value == "" {
		http.Error(rw, "metric value is empty", http.StatusBadRequest)
		return
	}

	err := counter.Store(name, value, h.storage)
	if err != nil {
		http.Error(rw, "metric value is bad", http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
}

func (h *handler) GetCounterHandler(rw http.ResponseWriter, r *http.Request) {
	name := strings.ToLower(chi.URLParam(r, "name"))
	log.Debug().
		Str("RequestURI", r.RequestURI).
		Str("name", name).
		Msg("")

	if name == "" {
		http.Error(rw, "metric name is empty", http.StatusNotFound)
		return
	}

	value, ok := counter.Get(name, h.storage)
	if !ok {
		http.Error(rw, "metric not found", http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := io.WriteString(rw, value)
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	rw.WriteHeader(http.StatusOK)
}

func (h *handler) PostGaugeHandler(rw http.ResponseWriter, r *http.Request) {
	name := strings.ToLower(chi.URLParam(r, "name"))
	log.Debug().
		Str("RequestURI", r.RequestURI).
		Str("name", name).
		Msg("")

	if name == "" {
		http.Error(rw, "metric name is empty", http.StatusNotFound)
		return
	}

	value := strings.ToLower(chi.URLParam(r, "value"))
	if value == "" {
		http.Error(rw, "metric value is empty", http.StatusBadRequest)
		return
	}

	err := gauge.Store(name, value, h.storage)
	if err != nil {
		http.Error(rw, "metric value is bad", http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
}

func (h *handler) GetGaugeHandler(rw http.ResponseWriter, r *http.Request) {
	name := strings.ToLower(chi.URLParam(r, "name"))
	log.Debug().
		Str("RequestURI", r.RequestURI).
		Str("name", name).
		Msg("")

	if name == "" {
		http.Error(rw, "metric name is empty", http.StatusNotFound)
		return
	}

	value, ok := gauge.Get(name, h.storage)
	if !ok {
		http.Error(rw, "metric not found", http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := io.WriteString(rw, value)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	rw.WriteHeader(http.StatusOK)
}
