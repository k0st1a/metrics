package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/k0st1a/metrics/internal/middleware/logging"
	"github.com/k0st1a/metrics/internal/models"
	"github.com/rs/zerolog/log"
)

const (
	badMetricType    = "metric type is bad"
	emptyMetricName  = "metric name is empty"
	emptyMetricValue = "metric value is empty"
	badMetricValue   = "metric value is bad"
	notFoundMetric   = "metric not found"

	badContentType = "bad Content-Type"
)

type storageService interface {
	Get(string) (string, bool)
	Store(string, string) error
}

type handler struct {
	storage storageService
}

type storageMetricService interface {
	GetGauge(string) (float64, bool)
	StoreGauge(string, float64)

	GetCounter(string) (int64, bool)
	StoreCounter(string, int64)
}

type handler2 struct {
	storage storageMetricService
}

func NewHandler2(s storageMetricService) *handler2 {
	return &handler2{
		storage: s,
	}
}

func NewHandler(s storageService) *handler {
	return &handler{
		storage: s,
	}
}

func BuildRouter(counter, gauge *handler, m *handler2) chi.Router {
	log.Debug().Msg("Make router")
	r := chi.NewRouter()

	r.Use(logging.Logging)

	log.Debug().Msg("the update handlers adding")
	r.Route("/update/", func(r chi.Router) {
		r.Post("/", m.PostUpdateHandler)
		r.Post("/counter/{name}/{value}", counter.PostMetricHandler)
		r.Post("/counter/", NotFoundHandler)
		r.Post("/gauge/{name}/{value}", gauge.PostMetricHandler)
		r.Post("/gauge/", NotFoundHandler)
	})
	log.Debug().Msg("the update handlers added")

	log.Debug().Msg("the value handlers adding")
	r.Route("/value/", func(r chi.Router) {
		r.Post("/", m.PostValueHandler)
		r.Get("/counter/{name}", counter.GetMetricHandler)
		r.Get("/gauge/{name}", gauge.GetMetricHandler)
	})
	log.Debug().Msg("the value handlers added")

	log.Debug().Msg("Custom NotFound handler adding")
	r.NotFound(BadRequestHandler)
	log.Debug().Msg("Custom NotFound handler added")

	return r
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

func (h *handler2) PostUpdateHandler(rw http.ResponseWriter, r *http.Request) {
	log.Info().
		Str("uri", r.RequestURI).
		Str("method", r.Method).
		Msg("")

	ct := r.Header.Get("Content-Type")
	if ct != "application/json" {
		http.Error(rw, badContentType, http.StatusBadRequest)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("io.ReadAll error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("")
		}
	}()

	m, err := models.Deserialize(b)
	if err != nil {
		log.Error().Err(err).Msg("models.Deserialize")
		http.Error(rw, "deserialize error", http.StatusBadRequest)
		return
	}

	switch m.MType {
	case "counter":
		v, ok := h.storage.GetCounter(m.ID)
		if ok {
			h.storage.StoreCounter(m.ID, *m.Delta+v)
		} else {
			h.storage.StoreCounter(m.ID, *m.Delta)
		}
	case "gauge":
		h.storage.StoreGauge(m.ID, *m.Value)
	default:
		http.Error(rw, badMetricType, http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (h *handler2) PostValueHandler(rw http.ResponseWriter, r *http.Request) {
	log.Info().
		Str("uri", r.RequestURI).
		Str("method", r.Method).
		Msg("")

	ct := r.Header.Get("Content-Type")
	if ct != "application/json" {
		http.Error(rw, badContentType, http.StatusBadRequest)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("rw.Body.Read")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("")
		}
	}()

	m, err := models.Deserialize(b)
	if err != nil {
		log.Error().Err(err).Msg("models.Deserialize")
		http.Error(rw, "deserialize error", http.StatusBadRequest)
		return
	}

	switch m.MType {
	case "counter":
		c, ok := h.storage.GetCounter(m.ID)
		if !ok {
			http.Error(rw, notFoundMetric, http.StatusNotFound)
			return
		}

		m.Delta = &c
	case "gauge":
		g, ok := h.storage.GetGauge(m.ID)
		if !ok {
			http.Error(rw, notFoundMetric, http.StatusNotFound)
			return
		}

		m.Value = &g
	default:
		http.Error(rw, badMetricType, http.StatusBadRequest)
		return
	}

	data2, err := models.Serialize(m)
	if err != nil {
		log.Error().Err(err).Msg("models.Serialize")
		http.Error(rw, "Serialize error", http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Agent-Type", "my-agent-type")

	_, err = rw.Write(data2)
	if err != nil {
		log.Error().Err(err).Msg("rw.Write error")
		return
	}
}
