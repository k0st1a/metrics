package json

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/k0st1a/metrics/internal/metrics"
	"github.com/k0st1a/metrics/internal/models"
	"github.com/rs/zerolog/log"
)

type reporter struct {
	url string
}

func NewReporter(a string) (*reporter, error) {
	url, err := url.JoinPath("http://", a, "/update/")
	if err != nil {
		return nil, fmt.Errorf("url.JoinPath error:%w", err)
	}

	return &reporter{
		url: url,
	}, nil
}

func (r reporter) DoReportsMetrics(c *http.Client, m *metrics.MyStats) {
	s := myStats2Metrics(m)
	for _, v := range s {
		r.doReportMetrics(c, &v)
	}
}

func (r reporter) doReportMetrics(c *http.Client, m *models.Metrics) {
	b, err := models.Serialize(m)
	if err != nil {
		log.Error().Err(err).Msg("models.Serialize")
		return
	}

	req, err := http.NewRequest(http.MethodPost, r.url, bytes.NewBuffer(b))
	if err != nil {
		log.Error().Err(err).Msg("http.NewRequest error")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("client do error")
		return
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("resp.Body.Close")
		}
	}()

	err = resp.Body.Close()
	if err != nil {
		log.Error().Err(err).Msg("resp.Body.Close error")
		return
	}
}

func myStats2Metrics(m *metrics.MyStats) []models.Metrics {
	return []models.Metrics{
		models.Metrics{
			ID:    "Alloc",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.Alloc),
		},
		models.Metrics{
			ID:    "BuckHashSys",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.BuckHashSys),
		},
		models.Metrics{
			ID:    "Frees",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.Frees),
		},
		models.Metrics{
			ID:    "GCSys",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.GCSys),
		},
		models.Metrics{
			ID:    "HeapAlloc",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.HeapAlloc),
		},
		models.Metrics{
			ID:    "HeapIdle",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.HeapIdle),
		},
		models.Metrics{
			ID:    "HeapInuse",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.HeapInuse),
		},
		models.Metrics{
			ID:    "HeapObjects",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.HeapObjects),
		},
		models.Metrics{
			ID:    "HeapReleased",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.HeapReleased),
		},
		models.Metrics{
			ID:    "HeapSys",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.HeapSys),
		},
		models.Metrics{
			ID:    "LastGC",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.LastGC),
		},
		models.Metrics{
			ID:    "Lookups",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.Lookups),
		},
		models.Metrics{
			ID:    "MCacheInuse",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.MCacheInuse),
		},
		models.Metrics{
			ID:    "MCacheSys",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.MCacheSys),
		},
		models.Metrics{
			ID:    "MSpanInuse",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.MSpanInuse),
		},
		models.Metrics{
			ID:    "MSpanSys",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.MSpanSys),
		},
		models.Metrics{
			ID:    "Mallocs",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.Mallocs),
		},
		models.Metrics{
			ID:    "NextGC",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.NextGC),
		},
		models.Metrics{
			ID:    "OtherSys",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.OtherSys),
		},
		models.Metrics{
			ID:    "PauseTotalNs",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.PauseTotalNs),
		},
		models.Metrics{
			ID:    "StackInuse",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.StackInuse),
		},
		models.Metrics{
			ID:    "StackSys",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.StackSys),
		},
		models.Metrics{
			ID:    "Sys",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.Sys),
		},
		models.Metrics{
			ID:    "TotalAlloc",
			MType: "gauge",
			Value: getAddressUInt64AsFloat64(m.MemStats.TotalAlloc),
		},
		models.Metrics{
			ID:    "NumForcedGC",
			MType: "gauge",
			Value: getAddressUInt32AsFloat64(m.MemStats.NumForcedGC),
		},
		models.Metrics{
			ID:    "NumGC",
			MType: "gauge",
			Value: getAddressUInt32AsFloat64(m.MemStats.NumGC),
		},
		models.Metrics{
			ID:    "GCCPUFraction",
			MType: "gauge",
			Value: &m.MemStats.GCCPUFraction,
		},
		models.Metrics{
			ID:    "PollCount",
			MType: "counter",
			Delta: getAddressUInt64AsInt64(m.PollCount),
		},
		models.Metrics{
			ID:    "RandomValue",
			MType: "gauge",
			Value: &m.RandomValue,
		},
	}
}

func getAddressUInt64AsInt64(v uint64) *int64 {
	v2 := int64(v)
	return &v2
}

func getAddressUInt64AsFloat64(v uint64) *float64 {
	v2 := float64(v)
	return &v2
}

func getAddressUInt32AsFloat64(v uint32) *float64 {
	v2 := float64(v)
	return &v2
}
