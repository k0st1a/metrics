package text

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/k0st1a/metrics/internal/handlers/json"
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

type Storage interface {
	GetGauge(ctx context.Context, name string) (*float64, error)
	StoreGauge(ctx context.Context, name string, value float64) error

	GetCounter(ctx context.Context, name string) (*int64, error)
	StoreCounter(ctx context.Context, name string, value int64) error

	GetAll(ctx context.Context) (gauge map[string]int64, counter map[string]float64, err error)
}

type handler struct {
	s Storage
}

func NewHandler(s Storage) *handler {
	return &handler{
		s: s,
	}
}

func BuildRouter(r *chi.Mux, h *handler) {
	r.Post("/update/{type}/{name}/{value}", h.PostMetricHandler)
	r.Post("/update/counter/", NotFoundHandler)
	r.Post("/update/gauge/", NotFoundHandler)

	r.Get("/", h.GetAllHandler)
	r.Get("/value/{type}/{name}", h.GetMetricHandler)
	r.Get("/value/{type}/{name}", h.GetMetricHandler)

	r.NotFound(BadRequestHandler)
}

func BadRequestHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
}

func NotFoundHandler(rw http.ResponseWriter, r *http.Request) {
	http.Error(rw, emptyMetricValue, http.StatusNotFound)
}

func (h *handler) GetAllHandler(rw http.ResponseWriter, r *http.Request) {
	const htmlTemplate = `Current metrics in form type/name/value:
{{range .}}{{.Type}}/{{.Name}}/{{.Value}}
{{end}}`

	ctx := r.Context()
	c, g, err := h.s.GetAll(ctx)
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

	ctx := r.Context()

	switch mtype {
	case "counter":
		c, err := str2counter(value)
		if err != nil {
			http.Error(rw, badMetricValue, http.StatusBadRequest)
			return
		}

		err = json.AddCounter(ctx, h.s, name, c)
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
		err = h.s.StoreGauge(ctx, name, g)
		if err != nil {
			log.Error().Err(err).Msg("storage gauge error")
			http.Error(rw, notFoundMetric, http.StatusInternalServerError)
			return
		}
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
}

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

	ctx := r.Context()

	switch mtype {
	case "counter":
		c, err := h.s.GetCounter(ctx, name)
		if err != nil {
			log.Error().Err(err).Msg("get counter error")
			http.Error(rw, notFoundMetric, http.StatusInternalServerError)
			return
		}
		if c == nil {
			http.Error(rw, notFoundMetric, http.StatusNotFound)
			return
		}
		value = counter2str(*c)
	case "gauge":
		g, err := h.s.GetGauge(ctx, name)
		if err != nil {
			log.Error().Err(err).Msg("get gauge error")
			http.Error(rw, notFoundMetric, http.StatusInternalServerError)
			return
		}
		if g == nil {
			http.Error(rw, notFoundMetric, http.StatusNotFound)
			return
		}
		value = gauge2str(*g)
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
