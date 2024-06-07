package json

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k0st1a/metrics/internal/handlers"
	"github.com/k0st1a/metrics/internal/pkg/hash"
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
		body               string
		contentType        string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "Metric type is bad",
			reqMethod:          http.MethodPost,
			reqPath:            "/value/",
			body:               `{"ID":"GaugeName","Type":"gauge"}`,
			contentType:        "application/json",
			expectedStatusCode: 400,
			expectedBody:       "metric type is bad\n",
		},
		{
			name:               "Get gauge metric with name GaugeName which not exists",
			reqMethod:          http.MethodPost,
			reqPath:            "/value/",
			body:               `{"id":"GaugeName","type":"gauge"}`,
			contentType:        "application/json",
			expectedStatusCode: 404,
			expectedBody:       "metric not found\n",
		},
		{
			name:               "Upload gauge metric with name GaugeName with value 123.3",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/",
			body:               `{"id":"GaugeName","type":"gauge","value":123.3}`,
			contentType:        "application/json",
			expectedStatusCode: 200,
			expectedBody:       "",
		},
		{
			name:               "Get gauge metric with name GaugeName with value 123.3",
			reqMethod:          http.MethodPost,
			reqPath:            "/value/",
			body:               `{"id":"GaugeName","type":"gauge"}`,
			contentType:        "application/json",
			expectedStatusCode: 200,
			expectedBody:       `{"value":123.3,"id":"GaugeName","type":"gauge"}`,
		},
		{
			name:               "Upload gauge metric with name GaugeName with value bad_value",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/",
			body:               `{"id":"GaugeName","type":"gauge","value":"bad_value"}`,
			contentType:        "application/json",
			expectedStatusCode: 400,
			expectedBody:       "deserialize error\n",
		},
		{
			name:               "Upload gauge metric with name GaugeName without value",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/",
			body:               `{"id":"GaugeName","type":"gauge"}`,
			contentType:        "application/json",
			expectedStatusCode: 400,
			expectedBody:       "metric value is nil\n",
		},
		{
			name:               "Upload gauge metric without name",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/",
			body:               `{"type":"gauge","value":123.3}`,
			contentType:        "application/json",
			expectedStatusCode: 400,
			expectedBody:       "metric id is empty\n",
		},
		{
			name:               "Upload unknown metric type",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/",
			body:               `{"id":"GaugeName","value":123.3}`,
			contentType:        "application/json",
			expectedStatusCode: 400,
			expectedBody:       "metric type is bad\n",
		},
		{
			name:               "Get counter metric with name CounterName which not exists",
			reqMethod:          http.MethodPost,
			reqPath:            "/value/",
			body:               `{"id":"CounterName","type":"counter"}`,
			contentType:        "application/json",
			expectedStatusCode: 404,
			expectedBody:       "metric not found\n",
		},
		{
			name:               "Upload counter metric with name CounterName with value 123",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/",
			body:               `{"id":"CounterName","type":"counter","delta":123}`,
			contentType:        "application/json",
			expectedStatusCode: 200,
			expectedBody:       "",
		},
		{
			name:               "Get counter metric with name CounterName with value 123",
			reqMethod:          http.MethodPost,
			reqPath:            "/value/",
			body:               `{"id":"CounterName","type":"counter"}`,
			contentType:        "application/json",
			expectedStatusCode: 200,
			expectedBody:       `{"delta":123,"id":"CounterName","type":"counter"}`,
		},
		{
			name:               "Upload counter metric with name CounterName with value bad_value",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/",
			body:               `{"id":"CounterName","type":"counter","value":"bad_value"}`,
			contentType:        "application/json",
			expectedStatusCode: 400,
			expectedBody:       "deserialize error\n",
		},
		{
			name:               "Upload counter metric with name CounterName without value",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/",
			body:               `{"id":"CounterName","type":"counter"}`,
			contentType:        "application/json",
			expectedStatusCode: 400,
			expectedBody:       "metric delta is nil\n",
		},
		{
			name:               "Upload counter metric without name",
			reqMethod:          http.MethodPost,
			reqPath:            "/update/",
			body:               `{"type":"counter","delta":123}`,
			contentType:        "application/json",
			expectedStatusCode: 400,
			expectedBody:       "metric id is empty\n",
		},
	}

	h := hash.New("")
	r := handlers.NewRouter(h)

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
			body := bytes.NewBuffer([]byte(test.body))
			req, err := http.NewRequest(test.reqMethod, testServer.URL+test.reqPath, body)
			if err != nil {
				t.Fatal(err)
			}
			if test.contentType != "" {
				req.Header.Set("Content-Type", test.contentType)
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
