package json

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

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
	mi := m.Metrics2MetricInfo()
	mms := []models.Metrics{}
	for _, v := range mi {
		mm, err := MetricInfo2Metrics(v)
		if err != nil {
			log.Error().Err(err).Msg("MetricInfo2Metrics error")
			continue
		}
		mms = append(mms, *mm)
	}
	return mms
}

func MetricInfo2Metrics(mi metrics.MetricInfo) (*models.Metrics, error) {
	switch mi.MType {
	case "gauge":
		v, err := str2float64(mi.Value)
		if err != nil {
			return nil, fmt.Errorf("str2float64 error:%w", err)
		}
		return &models.Metrics{
			ID:    mi.Name,
			MType: "gauge",
			Value: v,
		}, nil
	case "counter":
		v2, err := str2int64(mi.Value)
		if err != nil {
			return nil, fmt.Errorf("str2int64 error:%w", err)
		}
		return &models.Metrics{
			ID:    mi.Name,
			MType: "counter",
			Delta: v2,
		}, nil
	default:
		return nil, fmt.Errorf("unknown MType")
	}
}

func str2float64(s string) (*float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, fmt.Errorf("strconv.ParseFloat error:%w", err)
	}
	return &f, nil
}

func str2int64(s string) (*int64, error) {
	f, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("strconv.ParseInt error:%w", err)
	}
	return &f, nil
}
