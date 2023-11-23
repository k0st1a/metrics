package report

import (
	"net/http"
	"net/url"
	"time"

	"github.com/k0st1a/metrics/internal/metrics"
	"github.com/rs/zerolog/log"
)

func RunReportMetrics(addr string, client *http.Client, metrics *metrics.MyStats, reportInterval int) {
	ticker := time.NewTicker(time.Duration(reportInterval) * time.Second)

	for range ticker.C {
		reportMetrics(addr, client, metrics)
	}
}

func reportMetrics(addr string, client *http.Client, metrics *metrics.MyStats) {
	metrics.IncreasePollCount()
	prepared := metrics.Compose()
	for name, info := range prepared {
		reportMetric(addr, client, info.Type, name, info.Value)
	}
}

func reportMetric(addr string, client *http.Client, metricType, name, value string) {
	url, err := url.JoinPath("http://", addr, "/update/", metricType, "/", name, "/", value)
	if err != nil {
		log.Error().Err(err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		log.Error().Err(err)
		return
	}

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Content-Length", "0")

	response, err := client.Do(req)
	if err != nil {
		log.Error().Err(err)
		return
	}
	defer response.Body.Close()
}
