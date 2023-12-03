package json

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/k0st1a/metrics/internal/models"
	"github.com/rs/zerolog/log"
)

const (
	badMetricType  = "metric type is bad"
	notFoundMetric = "metric not found"

	badContentType = "bad Content-Type"
)

type storageMetricService interface {
	GetGauge(string) (float64, bool)
	StoreGauge(string, float64)

	GetCounter(string) (int64, bool)
	StoreCounter(string, int64)
}

type handler struct {
	storage storageMetricService
}

func NewHandler(s storageMetricService) *handler {
	return &handler{
		storage: s,
	}
}

func BuildRouter(r *chi.Mux, h *handler) {
	r.Post("/update/", h.PostUpdateHandler)
	r.Post("/value/", h.PostValueHandler)
}

func (h *handler) PostUpdateHandler(rw http.ResponseWriter, r *http.Request) {
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
		AddCounter(h.storage, m.ID, *m.Delta)
	case "gauge":
		h.storage.StoreGauge(m.ID, *m.Value)
	default:
		http.Error(rw, badMetricType, http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (h *handler) PostValueHandler(rw http.ResponseWriter, r *http.Request) {
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

	_, err = rw.Write(data2)
	if err != nil {
		log.Error().Err(err).Msg("rw.Write error")
		return
	}
}

func AddCounter(s storageMetricService, name string, value int64) {
	v, ok := s.GetCounter(name)
	if ok {
		s.StoreCounter(name, value+v)
	} else {
		s.StoreCounter(name, value)
	}
}
