package json

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/k0st1a/metrics/internal/agent/model"
	"github.com/k0st1a/metrics/internal/models"
	"github.com/rs/zerolog/log"
)

type report struct {
	client  *http.Client
	channel <-chan map[string]model.MetricInfoRaw
	address string
}

// NewReport - создание репортера, HTTP клиента, отправляющего метрики в формате JSON, где:
// * a - адрем сервера;
// * с - HTTP клиент;
// * ch - через данный канал получаем метрики для отправки на сервер.
func NewReport(a string, c *http.Client, ch <-chan map[string]model.MetricInfoRaw) *report {
	return &report{
		address: a,
		client:  c,
		channel: ch,
	}
}

// Do - запуск репортера.
func (r *report) Do() {
	for mi := range r.channel {
		log.Printf("mi:%v", mi)
		mi2 := model.RawMap2InfoList(mi)
		log.Printf("mi2:%v", mi2)
		ml := MetricsInfo2Metrics(mi2)
		log.Printf("ml:%v", ml)
		r.doReport(ml)
	}
}

func (r *report) doReport(m []models.Metrics) {
	b, err := models.SerializeList(m)
	if err != nil {
		log.Error().Err(err).Msg("models.SerializeList")
		return
	}

	url, err := url.JoinPath("http://", r.address, "/updates/")
	if err != nil {
		log.Error().Err(err).Msg("url.JoinPath")
		return
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		log.Error().Err(err).Msg("http.NewRequest error")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("client do error")
		return
	}

	err = resp.Body.Close()
	if err != nil {
		log.Error().Err(err).Msg("resp.Body.Close error")
		return
	}
}

// MetricsInfo2Metrics - преобразование списка метрики из "промежуточного" формата в "окончательный" формат
// для отправки на сервер.
func MetricsInfo2Metrics(mi []model.MetricInfo) []models.Metrics {
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

// MetricInfo2Metrics - преобразование метрики из "промежуточного" формате в "окончательный" формат
// для отправки на сервер.
func MetricInfo2Metrics(mi model.MetricInfo) (*models.Metrics, error) {
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
