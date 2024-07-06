// Package client is HTTP client for send metrics in JSON format.
package client

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	models "github.com/k0st1a/metrics/internal/adapters/api/http/model"
	"github.com/k0st1a/metrics/internal/agent/model"
	"github.com/rs/zerolog/log"
)

type client struct {
	client  *http.Client
	address string
}

func New(a string, t http.RoundTripper) *client {
	return &client{
		address: a,
		client: &http.Client{
			Transport: t,
		},
	}
}

func (c *client) Do(r model.MetricInfoRaw) error {
	m, err := Raw2Metric(r)
	if err != nil {
		return fmt.Errorf("Raw2Metric error:%w", err)
	}

	address, err := url.JoinPath("http://", c.address, "/update/")
	if err != nil {
		return fmt.Errorf("url join error:%w", err)
	}

	data, err := models.Serialize(m)
	if err != nil {
		return fmt.Errorf("models serialize error:%w", err)
	}

	req, err := http.NewRequest(http.MethodPost, address, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("http request error:%w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("client do error:%w", err)
	}

	log.Printf("Response StatusCode:%v", resp.StatusCode)

	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("response body close error:%w", err)
	}

	return nil
}

func (c *client) DoBatch(r []model.MetricInfoRaw) error {
	m := Raw2MetricList(r)

	address, err := url.JoinPath("http://", c.address, "/updates/")
	if err != nil {
		return fmt.Errorf("url join error:%w", err)
	}

	data, err := models.SerializeList(m)
	if err != nil {
		return fmt.Errorf("models serialize list error:%w", err)
	}

	req, err := http.NewRequest(http.MethodPost, address, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("http request error:%w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("client do error:%w", err)
	}

	log.Printf("Response StatusCode:%v", resp.StatusCode)

	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("response body close error:%w", err)
	}

	return nil
}

// Raw2MetricList - преобразование списка метрики из "сырого" формата в "окончательный" формат
// для отправки на сервер.
func Raw2MetricList(r []model.MetricInfoRaw) []models.Metrics {
	l := make([]models.Metrics, len(r))
	i := 0
	for _, v := range r {
		m, err := Raw2Metric(v)
		if err != nil {
			log.Error().Err(err).Msg("Raw2Metric error")
			continue
		}
		l[i] = *m
		i++
	}

	return l
}

// Raw2Metric - преобразование метрики из "сырого" формате в "окончательный" формат
// для отправки на сервер.
func Raw2Metric(r model.MetricInfoRaw) (*models.Metrics, error) {
	if r.Type == model.Gauge {
		v, ok := r.Value.(float64)
		if !ok {
			return nil, fmt.Errorf("for gauge(%+v) value type not float64", r)
		}

		return &models.Metrics{
			ID:    r.Name,
			MType: "gauge",
			Value: &v,
		}, nil
	}

	if r.Type == model.Counter {
		v, ok := r.Value.(int64)
		if !ok {
			return nil, fmt.Errorf("for counter(%+v) value type not int64", r)
		}

		return &models.Metrics{
			ID:    r.Name,
			MType: "counter",
			Delta: &v,
		}, nil
	}

	return nil, fmt.Errorf("unknown metric type(%+v)", r)
}
