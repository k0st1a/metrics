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

func BuildRouter() chi.Router {
	log.Debug().Msg("Make router")
	r := chi.NewRouter()

	log.Debug().Msg("POST handlers adding")
	r.Route("/update/", func(r chi.Router) {
		r.Post("/counter/{name}/{value}", PostCounterHandler)
		r.Post("/counter/", NotFoundHandler)
		r.Post("/gauge/{name}/{value}", PostGaugeHandler)
		r.Post("/gauge/", NotFoundHandler)
	})
	log.Debug().Msg("POST handlers added")

	log.Debug().Msg("GET handlers adding")
	r.Route("/value/", func(r chi.Router) {
		r.Get("/counter/{name}", GetCounterHandler)
		r.Get("/gauge/{name}", GetGaugeHandler)
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

	http.Error(rw, "metric value is empty", http.StatusNotFound)
}

func PostCounterHandler(rw http.ResponseWriter, r *http.Request) {
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

	err := counter.Store(name, value)
	if err != nil {
		http.Error(rw, "metric value is bad", http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
}

func GetCounterHandler(rw http.ResponseWriter, r *http.Request) {
	name := strings.ToLower(chi.URLParam(r, "name"))
	log.Debug().
		Str("RequestURI", r.RequestURI).
		Str("name", name).
		Msg("")

	if name == "" {
		http.Error(rw, "metric name is empty", http.StatusNotFound)
		return
	}

	value, ok := counter.Get(name)
	if !ok {
		http.Error(rw, "metric not found", http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(rw, value)
	rw.WriteHeader(http.StatusOK)
}

func PostGaugeHandler(rw http.ResponseWriter, r *http.Request) {
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

	err := gauge.Store(name, value)
	if err != nil {
		http.Error(rw, "metric value is bad", http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
}

func GetGaugeHandler(rw http.ResponseWriter, r *http.Request) {
	name := strings.ToLower(chi.URLParam(r, "name"))
	log.Debug().
		Str("RequestURI", r.RequestURI).
		Str("name", name).
		Msg("")

	if name == "" {
		http.Error(rw, "metric name is empty", http.StatusNotFound)
		return
	}

	value, ok := gauge.Get(name)
	if !ok {
		http.Error(rw, "metric not found", http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(rw, value)
	rw.WriteHeader(http.StatusOK)
}
