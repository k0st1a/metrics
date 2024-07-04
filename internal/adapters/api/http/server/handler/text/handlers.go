// Package text is HTTP text handler which work with DB via Storage interface.
package text

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/k0st1a/metrics/internal/pkg/retry"
	"github.com/k0st1a/metrics/internal/utils"
	"github.com/rs/zerolog/log"
)

type metricInfo struct {
	Type  string
	Name  string
	Value string
}

const (
	badMetricType    = "metric type is bad"
	emptyMetricName  = "metric name is empty"
	emptyMetricValue = "metric value is empty"
	badMetricValue   = "metric value is bad"
	notFoundMetric   = "metric not found"
)

// Storage - интерфейс работы с хранилищем метрик.
type Storage interface {
	// GetGauge - возвращает метрику типа gauge с именем name.
	GetGauge(ctx context.Context, name string) (*float64, error)
	// StoreGauge - сохраняет метрику типа gauge с именем name и значенем value.
	StoreGauge(ctx context.Context, name string, value float64) error

	// GetCounter - возвращает метрику типа gauge с именем name.
	GetCounter(ctx context.Context, name string) (*int64, error)
	// StoreCounter - сохраняет метрику типа counter с именем name и значенем value.
	StoreCounter(ctx context.Context, name string, value int64) error

	// StoreAll - сохраняет группу метрик типа counter и gauge.
	StoreAll(ctx context.Context, counter map[string]int64, gauge map[string]float64) error
	// GetAll - возвращает все метрики типа counter и gauge.
	GetAll(ctx context.Context) (counter map[string]int64, gauge map[string]float64, err error)
}

// Retryer - интерфейс повторного обращения к хранилищу.
type Retryer interface {
	Retry(ctx context.Context, check func(error) bool, fnc func() error) error
}

type handler struct {
	storage Storage
	retry   Retryer
}

// NewHandler - создание HTTP обработчика взаимодействия с хранилищем метрик.
func NewHandler(s Storage, r Retryer) *handler {
	return &handler{
		storage: s,
		retry:   r,
	}
}

// BuildRouter - формирование маршрута для HTTP обработчика.
func BuildRouter(r *chi.Mux, h *handler) {
	r.Post("/update/{type}/{name}/{value}", h.PostMetricHandler)
	r.Post("/update/counter/", NotFoundHandler)
	r.Post("/update/gauge/", NotFoundHandler)

	r.Get("/", h.GetAllHandler)
	r.Get("/value/{type}/{name}", h.GetMetricHandler)

	r.NotFound(BadRequestHandler)
}

// BadRequestHandler - обработчик когда не найден ни один из путей.
func BadRequestHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
}

// NotFoundHandler - обработчик когда в пути не задано имя метрики.
func NotFoundHandler(rw http.ResponseWriter, r *http.Request) {
	http.Error(rw, emptyMetricValue, http.StatusNotFound)
}

// GetAllHandler - обработчик для получения всех метрик.
func (h *handler) GetAllHandler(rw http.ResponseWriter, r *http.Request) {
	const htmlTemplate = `Current metrics in form type/name/value:
{{range .}}{{.Type}}/{{.Name}}/{{.Value}}
{{end}}`

	var (
		c   map[string]int64
		g   map[string]float64
		err error
	)

	err = h.retry.Retry(r.Context(), retry.IsConnectionException, func() error {
		c, g, err = h.storage.GetAll(r.Context())
		//nolint // Не за чем оборачивать ошибку
		return err
	})
	if err != nil {
		log.Error().Err(err).Msg("get metrics error")
		return
	}

	m := make([]metricInfo, 0)

	for n, v := range c {
		m = append(m, metricInfo{Type: "counter", Name: n, Value: counter2str(v)})
	}

	for n, v := range g {
		m = append(m, metricInfo{Type: "gauge", Name: n, Value: gauge2str(v)})
	}

	t := template.New("myTemplate")

	t, err = t.Parse(htmlTemplate)
	if err != nil {
		log.Error().Err(err).Msg("t.Parse error")
		return
	}

	rw.Header().Set("Content-Type", "text/html")

	err = t.Execute(rw, m)
	if err != nil {
		log.Error().Err(err).Msg("t.Execute error")
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// PostMetricHandler - обработчик для сохранения метрики.
func (h *handler) PostMetricHandler(rw http.ResponseWriter, r *http.Request) {
	mtype := strings.ToLower(chi.URLParam(r, "type"))
	if !checkType(mtype) {
		http.Error(rw, badMetricType, http.StatusBadRequest)
		return
	}

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

	switch mtype {
	case "counter":
		c, err := str2counter(value)
		if err != nil {
			http.Error(rw, badMetricValue, http.StatusBadRequest)
			return
		}

		err = h.retry.Retry(r.Context(), retry.IsConnectionException, func() error {
			//nolint // Не за чем оборачивать ошибку
			return h.storage.StoreCounter(r.Context(), name, c)
		})
		if err != nil {
			log.Error().Err(err).Msg("add counter error")
			http.Error(rw, notFoundMetric, http.StatusInternalServerError)
			return
		}
	case "gauge":
		g, err := str2gauge(value)
		if err != nil {
			http.Error(rw, badMetricValue, http.StatusBadRequest)
			return
		}
		err = h.retry.Retry(r.Context(), retry.IsConnectionException, func() error {
			//nolint // Не за чем оборачивать ошибку
			return h.storage.StoreGauge(r.Context(), name, g)
		})
		if err != nil {
			log.Error().Err(err).Msg("storage gauge error")
			http.Error(rw, notFoundMetric, http.StatusInternalServerError)
			return
		}
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
}

// GetMetricHandler - обработчик для получения метрики.
func (h *handler) GetMetricHandler(rw http.ResponseWriter, r *http.Request) {
	mtype := strings.ToLower(chi.URLParam(r, "type"))
	if !checkType(mtype) {
		http.Error(rw, badMetricType, http.StatusBadRequest)
		return
	}

	name := strings.ToLower(chi.URLParam(r, "name"))
	if name == "" {
		http.Error(rw, emptyMetricName, http.StatusNotFound)
		return
	}

	var value string

	switch mtype {
	case "counter":
		var (
			c   *int64
			err error
		)
		err = h.retry.Retry(r.Context(), retry.IsConnectionException, func() error {
			c, err = h.storage.GetCounter(r.Context(), name)
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
			if c != nil {
				value = counter2str(*c)
			} else {
				log.Error().Err(err).Msg("nil counter value")
				http.Error(rw, notFoundMetric, http.StatusInternalServerError)
			}
		}
	case "gauge":
		var (
			g   *float64
			err error
		)
		err = h.retry.Retry(r.Context(), retry.IsConnectionException, func() error {
			g, err = h.storage.GetGauge(r.Context(), name)
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
			if g != nil {
				value = gauge2str(*g)
			} else {
				log.Error().Err(err).Msg("nil gauge value")
				http.Error(rw, notFoundMetric, http.StatusInternalServerError)
			}
		}
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")

	_, err := rw.Write([]byte(value))
	if err != nil {
		log.Error().Err(err).Msg("rw.Write error")
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func checkType(t string) bool {
	switch t {
	case "counter":
		return true
	case "gauge":
		return true
	default:
		return false
	}
}

func str2counter(s string) (int64, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return v, fmt.Errorf("parse int error:%w", err)
	}

	return v, nil
}

func counter2str(i int64) string {
	return strconv.FormatInt(i, 10)
}

func str2gauge(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return v, fmt.Errorf("parse float error:%w", err)
	}

	return v, nil
}

func gauge2str(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
