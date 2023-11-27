package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

const (
	emptyMetricName  = "metric name is empty"
	emptyMetricValue = "metric value is empty"
	badMetricValue   = "metric value is bad"
	notFoundMetric   = "metric not found"
)

type storageService interface {
	Get(string) (string, bool)
	Store(string, string) error
}

type handler struct {
	storage storageService
}

func NewHandler(s storageService) *handler {
	return &handler{
		storage: s,
	}
}

func BuildRouter(counter, gauge *handler) chi.Router {
	log.Debug().Msg("Make router")
	r := chi.NewRouter()

	log.Debug().Msg("POST handlers adding")
	r.Route("/update/", func(r chi.Router) {
		r.Post("/counter/{name}/{value}", counter.PostMetricHandler)
		r.Post("/counter/", NotFoundHandler)
		r.Post("/gauge/{name}/{value}", gauge.PostMetricHandler)
		r.Post("/gauge/", NotFoundHandler)
	})
	log.Debug().Msg("POST handlers added")

	log.Debug().Msg("GET handlers adding")
	r.Route("/value/", func(r chi.Router) {
		r.Get("/counter/{name}", counter.GetMetricHandler)
		r.Get("/gauge/{name}", gauge.GetMetricHandler)
	})
	log.Debug().Msg("GET handlers added")

	log.Debug().Msg("Custom NotFound handler adding")
	r.NotFound(BadRequestHandler)
	log.Debug().Msg("Custom NotFound handler added")

	return r
}

func BadRequestHandler(rw http.ResponseWriter, r *http.Request) {
	log.Debug().
		Str("RequestURI", r.RequestURI).
		Msg("")

	rw.WriteHeader(http.StatusBadRequest)
}

func NotFoundHandler(rw http.ResponseWriter, r *http.Request) {
	log.Debug().
		Str("RequestURI", r.RequestURI).
		Msg("")

	http.Error(rw, emptyMetricValue, http.StatusNotFound)
}

func (h *handler) PostMetricHandler(rw http.ResponseWriter, r *http.Request) {
	name := strings.ToLower(chi.URLParam(r, "name"))
	log.Debug().
		Str("RequestURI", r.RequestURI).
		Str("name", name).
		Msg("")

	if name == "" {
		http.Error(rw, emptyMetricName, http.StatusNotFound)
		return
	}

	value := strings.ToLower(chi.URLParam(r, "value"))
	if value == "" {
		http.Error(rw, emptyMetricValue, http.StatusBadRequest)
		return
	}

	err := h.storage.Store(name, value)
	if err != nil {
		http.Error(rw, badMetricValue, http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
}

func (h *handler) GetMetricHandler(rw http.ResponseWriter, r *http.Request) {
	name := strings.ToLower(chi.URLParam(r, "name"))
	log.Debug().
		Str("RequestURI", r.RequestURI).
		Str("name", name).
		Msg("")

	if name == "" {
		http.Error(rw, emptyMetricName, http.StatusNotFound)
		return
	}

	value, ok := h.storage.Get(name)
	if !ok {
		http.Error(rw, notFoundMetric, http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := rw.Write([]byte(value))
	if err != nil {
		log.Error().Err(err).Msg("rw.Write error")
		return
	}

	rw.WriteHeader(http.StatusOK)
}
