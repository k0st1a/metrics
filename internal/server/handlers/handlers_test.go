package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k0st1a/metrics/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCounter(t *testing.T) {

	fmt.Println("Running logger")
	logger.Run()
	defer logger.Close()

	type want struct {
		contentType string
		statusCode  int
		body        string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "check update gauge metric with name GaugeName with value 123.3",
			request: "/update/gauge/GaugeName/123.3",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  200,
				body:        "",
			},
		},
		{
			name:    "check update gauge metric with name GaugeName with value bad_value",
			request: "/update/gauge/GaugeName/bad_value",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				body:        "gauge value is bad\n",
			},
		},
		{
			name:    "check update gauge metric with name GaugeName without value",
			request: "/update/gauge/GaugeName/",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				body:        "gauge value is empty\n",
			},
		},
		{
			name:    "check update gauge metric without name",
			request: "/update/gauge/",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
				body:        "gauge name is empty\n",
			},
		},
		{
			name:    "check bad gauge request",
			request: "/update/gauges",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				body:        "bad gauge request\n",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.request, nil)
			w := httptest.NewRecorder()
			Gauge(w, request)

			res := w.Result()
			assert.Equal(t, test.want.statusCode, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.body, string(resBody))
		})
	}
}
