package text

import (
	"net/http"
	"net/url"

	"github.com/k0st1a/metrics/internal/metrics"
	"github.com/rs/zerolog/log"
)

type reporter struct {
	addr string
}

func NewReporter(a string) (*reporter, error) {
	return &reporter{
		addr: a,
	}, nil
}

func (r reporter) DoReportsMetrics(c *http.Client, m *metrics.MyStats) {
	s := m.Metrics2MetricInfo()
	for _, v := range s {
		r.doReportMetrics(c, v)
	}
}

func (r reporter) doReportMetrics(c *http.Client, m metrics.MetricInfo) {
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

	response, err := c.Do(req)
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
