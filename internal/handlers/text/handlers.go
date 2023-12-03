package text

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

const (
	badMetricType    = "metric type is bad"
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

func BuildRouter(r *chi.Mux, counter, gauge *handler) {
	r.Post("/update/counter/{name}/{value}", counter.PostMetricHandler)
	r.Post("/update/counter/", NotFoundHandler)
	r.Post("/update/gauge/{name}/{value}", gauge.PostMetricHandler)
	r.Post("/update/gauge/", NotFoundHandler)

	r.Get("/value/counter/{name}", counter.GetMetricHandler)
	r.Get("/value/gauge/{name}", gauge.GetMetricHandler)

	r.NotFound(BadRequestHandler)
}

func BadRequestHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
}

func NotFoundHandler(rw http.ResponseWriter, r *http.Request) {
	http.Error(rw, emptyMetricValue, http.StatusNotFound)
}

func (h *handler) PostMetricHandler(rw http.ResponseWriter, r *http.Request) {
	name := strings.ToLower(chi.URLParam(r, "name"))

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
