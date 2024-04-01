package text

import (
	"net/http"
	"net/url"

	"github.com/k0st1a/metrics/internal/agent/model"
	"github.com/rs/zerolog/log"
)

type Doer interface {
	Do()
}

type Metrics2MetricInfoer interface {
	Metrics2MetricInfo() []model.MetricInfo
}

type report struct {
	c    *http.Client
	m    Metrics2MetricInfoer
	addr string
}

func NewReport(a string, c *http.Client, m Metrics2MetricInfoer) Doer {
	return &report{
		addr: a,
		c:    c,
		m:    m,
	}
}

func (r report) Do() {
	s := r.m.Metrics2MetricInfo()
	for _, v := range s {
		r.doReport(v)
	}
}

func (r report) doReport(m model.MetricInfo) {
	url, err := url.JoinPath("http://", r.addr, "/update/", m.MType, "/", m.Name, "/", m.Value)
	if err != nil {
		log.Error().Err(err).Msg("url.JoinPath error")
		return
	}

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		log.Error().Err(err).Msg("http.NewRequest error")
		return
	}

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Content-Length", "0")

	response, err := r.c.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("client do error")
		return
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("response.Body.Close")
		}
	}()
}
