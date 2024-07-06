// Package text contains the reporter code which work with text format.
package text

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

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
	var (
		value string
		mtype string
	)

	switch r.Type {
	case model.Gauge:
		v, ok := r.Value.(float64)
		if !ok {
			return fmt.Errorf("for gauge(%v) value type not float64", r)
		}

		value = strconv.FormatFloat(v, 'g', -1, 64)
		mtype = "gauge"
	case model.Counter:
		v, ok := r.Value.(int64)
		if !ok {
			return fmt.Errorf("for counter(%v) value type not int64", r)
		}

		value = strconv.FormatInt(v, 10)
		mtype = "counter"
	default:
		return fmt.Errorf(
			"bad metric value type. type:%v, name:%v, actual type:%v",
			r.Type, r.Name, reflect.TypeOf(r.Value))
	}

	url, err := url.JoinPath("http://", c.address, "/update/", mtype, "/", r.Name, "/", value)
	if err != nil {
		return fmt.Errorf("url join error:%w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("http request error:%w", err)
	}

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Content-Length", "0")

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
