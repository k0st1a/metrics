package text

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/k0st1a/metrics/internal/metrics"
	"github.com/rs/zerolog/log"
)

type textReporter struct {
	addr string
}

func NewTextReporter(a string) (*textReporter, error) {
	return &textReporter{
		addr: a,
	}, nil
}

func (r textReporter) DoReportsMetrics(c *http.Client, m *metrics.MyStats) {
	s := myStats2metricsInfo(m)
	for _, v := range s {
		r.doReportMetrics(c, v)
	}
}

func (r textReporter) doReportMetrics(c *http.Client, m metricInfo) {
	url, err := url.JoinPath("http://", r.addr, "/update/", m.mtype, "/", m.name, "/", m.value)
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

type metricInfo struct {
	name  string
	mtype string
	value string
}

func myStats2metricsInfo(m *metrics.MyStats) []metricInfo {
	return []metricInfo{
		metricInfo{
			name:  "Alloc",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.Alloc, 10),
		},
		metricInfo{
			name:  "BuckHashSys",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.BuckHashSys, 10),
		},
		metricInfo{
			name:  "Frees",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.Frees, 10),
		},
		metricInfo{
			name:  "GCSys",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.GCSys, 10),
		},
		metricInfo{
			name:  "HeapAlloc",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.HeapAlloc, 10),
		},
		metricInfo{
			name:  "HeapIdle",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.HeapIdle, 10),
		},
		metricInfo{
			name:  "HeapInuse",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.HeapInuse, 10),
		},
		metricInfo{
			name:  "HeapObjects",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.HeapObjects, 10),
		},
		metricInfo{
			name:  "HeapReleased",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.HeapReleased, 10),
		},
		metricInfo{
			name:  "HeapSys",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.HeapSys, 10),
		},
		metricInfo{
			name:  "LastGC",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.LastGC, 10),
		},
		metricInfo{
			name:  "Lookups",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.Lookups, 10),
		},
		metricInfo{
			name:  "MCacheInuse",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.MCacheInuse, 10),
		},
		metricInfo{
			name:  "MCacheSys",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.MCacheSys, 10),
		},
		metricInfo{
			name:  "MSpanInuse",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.MSpanInuse, 10),
		},
		metricInfo{
			name:  "MSpanSys",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.MSpanSys, 10),
		},
		metricInfo{
			name:  "Mallocs",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.Mallocs, 10),
		},
		metricInfo{
			name:  "NextGC",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.NextGC, 10),
		},
		metricInfo{
			name:  "OtherSys",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.OtherSys, 10),
		},
		metricInfo{
			name:  "PauseTotalNs",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.PauseTotalNs, 10),
		},
		metricInfo{
			name:  "StackInuse",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.StackInuse, 10),
		},
		metricInfo{
			name:  "StackSys",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.StackSys, 10),
		},
		metricInfo{
			name:  "Sys",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.Sys, 10),
		},
		metricInfo{
			name:  "TotalAlloc",
			mtype: "gauge",
			value: strconv.FormatUint(m.MemStats.TotalAlloc, 10),
		},
		metricInfo{
			name:  "NumForcedGC",
			mtype: "gauge",
			value: strconv.FormatUint(uint64(m.MemStats.NumForcedGC), 10),
		},
		metricInfo{
			name:  "NumGC",
			mtype: "gauge",
			value: strconv.FormatUint(uint64(m.MemStats.NumGC), 10),
		},
		metricInfo{
			name:  "GCCPUFraction",
			mtype: "gauge",
			value: strconv.FormatFloat(m.MemStats.GCCPUFraction, 'g', -1, 64),
		},
		metricInfo{
			name:  "PollCount",
			mtype: "counter",
			value: strconv.FormatUint(m.PollCount, 10),
		},
		metricInfo{
			name:  "RandomValue",
			mtype: "gauge",
			value: strconv.FormatFloat(m.RandomValue, 'g', -1, 64),
		},
	}
}
