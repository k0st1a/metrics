package json

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/k0st1a/metrics/internal/models"
	"github.com/k0st1a/metrics/internal/pkg/retry"
	"github.com/k0st1a/metrics/internal/utils"
	"github.com/rs/zerolog/log"
)

const (
	badMetricType  = "metric type is bad"
	notFoundMetric = "metric not found"
)

type Storage interface {
	GetGauge(ctx context.Context, name string) (*float64, error)
	StoreGauge(ctx context.Context, name string, value float64) error

	GetCounter(ctx context.Context, name string) (*int64, error)
	StoreCounter(ctx context.Context, name string, value int64) error

	StoreAll(ctx context.Context, counter map[string]int64, gauge map[string]float64) error
	GetAll(ctx context.Context) (counter map[string]int64, gauge map[string]float64, err error)
}

type Retryer interface {
	Retry(ctx context.Context, check func(error) bool, fnc func() error) error
}

type handler struct {
	storage Storage
	retry   Retryer
}

func NewHandler(s Storage, r Retryer) *handler {
	return &handler{
		storage: s,
		retry:   r,
	}
}

func BuildRouter(r *chi.Mux, h *handler) {
	r.With(contentType).Post("/updates/", h.PostUpdatesHandler)
	r.With(contentType).Post("/update/", h.PostUpdateHandler)
	r.With(contentType).Post("/value/", h.PostValueHandler)
}

func (h *handler) PostUpdatesHandler(rw http.ResponseWriter, r *http.Request) {
	log.Info().
		Str("uri", r.RequestURI).
		Str("method", r.Method).
		Msg("")

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("io.ReadAll error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	m, err := models.DeserializeList(b)
	if err != nil {
		log.Error().Err(err).Msg("models.Deserialize error")
		http.Error(rw, "deserialize error", http.StatusBadRequest)
		return
	}

	g := make(map[string]float64)
	c := make(map[string]int64)

	for _, v := range m {
		switch v.MType {
		case "counter":
			cv, ok := c[v.ID]
			switch ok {
			case true:
				c[v.ID] = cv + *v.Delta
			default:
				c[v.ID] = *v.Delta
			}
		case "gauge":
			g[v.ID] = *v.Value
		default:
			log.Error().
				Str("unknown MType", v.MType).
				Msg("")
		}
	}

	log.Printf("Store\nCounters:%+v\nGauges:%+v\n", c, g)

	err = h.retry.Retry(r.Context(), retry.IsConnectionException, func() error {
		//nolint // Не за чем оборачивать ошибку
		return h.storage.StoreAll(r.Context(), c, g)
	})
	if err != nil {
		log.Error().Err(err).Msg("s.StoreAll error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (h *handler) PostUpdateHandler(rw http.ResponseWriter, r *http.Request) {
	log.Info().
		Str("uri", r.RequestURI).
		Str("method", r.Method).
		Msg("")

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("io.ReadAll error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	m, err := models.Deserialize(b)
	if err != nil {
		log.Error().Err(err).Msg("models.Deserialize")
		http.Error(rw, "deserialize error", http.StatusBadRequest)
		return
	}

	switch m.MType {
	case "counter":
		log.Printf("Post Update counter, name(%v), value(%v)", m.ID, *m.Delta)
		err = h.retry.Retry(r.Context(), retry.IsConnectionException, func() error {
			//nolint // Не за чем оборачивать ошибку
			return h.storage.StoreCounter(r.Context(), m.ID, *m.Delta)
		})
		if err != nil {
			log.Error().Err(err).Msg("h.storage.StoreCounter error")
			http.Error(rw, "store counter error", http.StatusInternalServerError)
			return
		}
	case "gauge":
		log.Printf("Post Update gauge, name(%v), value(%v)", m.ID, *m.Value)
		err = h.retry.Retry(r.Context(), retry.IsConnectionException, func() error {
			//nolint // Не за чем оборачивать ошибку
			return h.storage.StoreGauge(r.Context(), m.ID, *m.Value)
		})
		if err != nil {
			log.Error().Err(err).Msg("h.storage.StorageGauge error")
			http.Error(rw, "storage gauge error", http.StatusInternalServerError)
			return
		}
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

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("rw.Body.Read")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	m, err := models.Deserialize(b)
	if err != nil {
		log.Error().Err(err).Msg("models.Deserialize")
		http.Error(rw, "deserialize error", http.StatusBadRequest)
		return
	}

	switch m.MType {
	case "counter":
		var c *int64
		err = h.retry.Retry(r.Context(), retry.IsConnectionException, func() error {
			c, err = h.storage.GetCounter(r.Context(), m.ID)
			//nolint // Не за чем оборачивать ошибку
			return err
		})
		switch {
		case errors.Is(err, utils.ErrMetricsNoCounter):
			http.Error(rw, notFoundMetric, http.StatusNotFound)
			return
		case err != nil:
			log.Error().Err(err).Msg("get counter error")
			http.Error(rw, notFoundMetric, http.StatusInternalServerError)
			return
		default:
			m.Delta = c
		}
	case "gauge":
		var g *float64
		err = h.retry.Retry(r.Context(), retry.IsConnectionException, func() error {
			g, err = h.storage.GetGauge(r.Context(), m.ID)
			//nolint // Не за чем оборачивать ошибку
			return err
		})
		switch {
		case errors.Is(err, utils.ErrMetricsNoGauge):
			http.Error(rw, notFoundMetric, http.StatusNotFound)
			return
		case err != nil:
			log.Error().Err(err).Msg("get gauge error")
			http.Error(rw, notFoundMetric, http.StatusInternalServerError)
			return
		default:
			m.Value = g
		}
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

func contentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(rw, "bad Content-Type", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(rw, r)
	})
}
