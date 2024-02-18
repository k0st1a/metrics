package json

import (
	"context"
	"fmt"
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

type Storage interface {
	GetGauge(ctx context.Context, name string) (*float64, error)
	StoreGauge(ctx context.Context, name string, value float64) error

	GetCounter(ctx context.Context, name string) (*int64, error)
	StoreCounter(ctx context.Context, name string, value int64) error
}

type handler struct {
	storage Storage
}

func NewHandler(s Storage) *handler {
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

	ctx := r.Context()

	switch m.MType {
	case "counter":
		log.Printf("Post Update counter, name(%v), value(%v)", m.ID, *m.Delta)
		AddCounter(ctx, h.storage, m.ID, *m.Delta)
	case "gauge":
		log.Printf("Post Update gauge, name(%v), value(%v)", m.ID, *m.Value)
		h.storage.StoreGauge(ctx, m.ID, *m.Value)
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

	ctx := r.Context()

	switch m.MType {
	case "counter":
		c, err := h.storage.GetCounter(ctx, m.ID)
		if err != nil {
			log.Error().Err(err).Msg("get counter error")
			http.Error(rw, notFoundMetric, http.StatusInternalServerError)
			return
		}
		if c == nil {
			http.Error(rw, notFoundMetric, http.StatusNotFound)
			return
		}
		m.Delta = c
	case "gauge":
		g, err := h.storage.GetGauge(ctx, m.ID)
		if err != nil {
			log.Error().Err(err).Msg("get gauge error")
			http.Error(rw, notFoundMetric, http.StatusInternalServerError)
			return
		}
		if g == nil {
			http.Error(rw, notFoundMetric, http.StatusNotFound)
			return
		}
		m.Value = g
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

func AddCounter(ctx context.Context, s Storage, name string, value int64) error {
	v, err := s.GetCounter(ctx, name)
	if err != nil {
		return fmt.Errorf("get counter error:%w", err)
	}

	v2 := value
	if v != nil {
		v2 += (*v)
	}

	err = s.StoreCounter(ctx, name, v2)
	if err != nil {
		return fmt.Errorf("store counter error:%w", err)
	}

	return nil
}
