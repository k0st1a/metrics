package text

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k0st1a/metrics/internal/handlers"
	"github.com/k0st1a/metrics/internal/pkg/retry"
	"github.com/k0st1a/metrics/internal/storage/inmemory"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricHandler(t *testing.T) {
	tests := []struct {
		name               string
		reqMethod          string
		reqPath            string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "check get gauge metric with name GaugeName which not exists",
			reqMethod:          http.MethodGet,
			reqPath:            "/value/gauge/GaugeName",
			expectedStatusCode: 404,
			expectedBody:       "metric not found\n",
		},
		{
			name:               "check update gauge metric with name GaugeName with value 123.3",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/gauge/GaugeName/123.3",
			expectedStatusCode: 200,
			expectedBody:       "",
		},
		{
			name:               "check get gauge metric with name GaugeName with value 123.3",
			reqMethod:          http.MethodGet,
			reqPath:            "/value/gauge/GaugeName",
			expectedStatusCode: 200,
			expectedBody:       "123.3",
		},
		{
			name:               "check update gauge metric with name GaugeName with value bad_value",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/gauge/GaugeName/bad_value",
			expectedStatusCode: 400,
			expectedBody:       "metric value is bad\n",
		},
		{
			name:               "check update gauge metric with name GaugeName without value",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/gauge/GaugeName/",
			expectedStatusCode: 400,
			expectedBody:       "",
		},
		{
			name:               "check update gauge metric without name",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/gauge/",
			expectedStatusCode: 404,
			expectedBody:       "metric value is empty\n",
		},
		{
			name:               "check bad gauge request",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/gauges",
			expectedStatusCode: 400,
			expectedBody:       "",
		},
		{
			name:               "check unknown metric type",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/unknown/testCounter/100",
			expectedStatusCode: 400,
			expectedBody:       "metric type is bad\n",
		},
		{
			name:               "check get all metrics",
			reqMethod:          http.MethodGet,
			reqPath:            "/",
			expectedStatusCode: 200,
			expectedBody:       "Current metrics in form type/name/value:\ngauge/gaugename/123.3\n",
		},
		{
			name:               "check get counter metric with name CounterName which not exists",
			reqMethod:          http.MethodGet,
			reqPath:            "/value/counter/CounterName",
			expectedStatusCode: 404,
			expectedBody:       "metric not found\n",
		},
		{
			name:               "check update counter metric with name CounterName with value 123",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/counter/CounterName/123",
			expectedStatusCode: 200,
			expectedBody:       "",
		},
		{
			name:               "check get counter metric with name CounterName with value 123",
			reqMethod:          http.MethodGet,
			reqPath:            "/value/counter/CounterName",
			expectedStatusCode: 200,
			expectedBody:       "123",
		},
		{
			name:               "check update counter metric with name CounterName with value bad_value",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/counter/CounterName/bad_value",
			expectedStatusCode: 400,
			expectedBody:       "metric value is bad\n",
		},
		{
			name:               "check update counter metric with name CounterName without value",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/counter/CounterName/",
			expectedStatusCode: 400,
			expectedBody:       "",
		},
		{
			name:               "check update counter metric without name",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/counter/",
			expectedStatusCode: 404,
			expectedBody:       "metric value is empty\n",
		},
		{
			name:               "check bad counter request",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/counters",
			expectedStatusCode: 400,
			expectedBody:       "",
		},
		{
			name:               "check get all metrics",
			reqMethod:          http.MethodGet,
			reqPath:            "/",
			expectedStatusCode: 200,
			expectedBody:       "Current metrics in form type/name/value:\ncounter/countername/123\ngauge/gaugename/123.3\n",
		},
	}

	r := handlers.NewRouter(nil)
	s := inmemory.NewStorage()
	rt := retry.New()
	th := NewHandler(s, rt)

	BuildRouter(r, th)

	testServer := httptest.NewServer(r)
	defer testServer.Close()

	testClient := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.reqMethod, testServer.URL+test.reqPath, nil)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := testClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			err = resp.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, test.expectedBody, string(respBody))
			require.Equal(t, test.expectedStatusCode, resp.StatusCode)
		})
	}
}
