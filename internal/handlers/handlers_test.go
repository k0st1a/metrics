package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k0st1a/metrics/internal/storage"
	"github.com/k0st1a/metrics/internal/storage/counter"
	"github.com/k0st1a/metrics/internal/storage/gauge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostMetricHandler(t *testing.T) {
	tests := []struct {
		name               string
		reqMethod          string
		reqPath            string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "check update gauge metric with name GaugeName with value 123.3",
			reqPath:            "/update/gauge/GaugeName/123.3",
			expectedStatusCode: 200,
			expectedBody:       "",
		},
		{
			name:               "check update gauge metric with name GaugeName with value bad_value",
			reqPath:            "/update/gauge/GaugeName/bad_value",
			expectedStatusCode: 400,
			expectedBody:       "metric value is bad\n",
		},
		{
			name:               "check update gauge metric with name GaugeName without value",
			reqPath:            "/update/gauge/GaugeName/",
			expectedStatusCode: 400,
			expectedBody:       "",
		},
		{
			name:               "check update gauge metric without name",
			reqPath:            "/update/gauge/",
			expectedStatusCode: 404,
			expectedBody:       "metric value is empty\n",
		},
		{
			name:               "check bad gauge request",
			reqPath:            "/update/gauges",
			expectedStatusCode: 400,
			expectedBody:       "",
		},
		{
			name:               "check unknown metric type",
			reqPath:            "/update/unknown/testCounter/100",
			expectedStatusCode: 400,
			expectedBody:       "",
		},
	}

	s := storage.NewStorage()
	gs := gauge.NewGaugeStorage(s)
	cs := counter.NewCounterStorage(s)
	csh := NewHandler(cs)
	gsh := NewHandler(gs)
	mh := NewHandler2(s)

	testServer := httptest.NewServer(BuildRouter(csh, gsh, mh))
	defer testServer.Close()

	testClient := &http.Client{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, testServer.URL+test.reqPath, nil)
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
