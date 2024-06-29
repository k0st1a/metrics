package trustedsubnet

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/3th1nk/cidr"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrustedSubnetError(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		subnet         string
		xRealIP        string
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name:           "check request ip is not trusted",
			body:           "some body",
			subnet:         "192.168.1.0/24",
			xRealIP:        "10.25.88.22",
			expectedStatus: http.StatusForbidden,
			expectedBody:   []byte("request ip is not trusted\n"),
		},
		{
			name:           "check request ip is trusted",
			body:           "some body",
			subnet:         "192.168.1.0/24",
			xRealIP:        "192.168.1.1",
			expectedStatus: http.StatusOK,
			expectedBody:   []byte("some body"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cidr, err := cidr.Parse(test.subnet)
			require.NoError(t, err)

			r := chi.NewRouter()
			r.Use(New(cidr))
			respRecorder := httptest.NewRecorder()

			r.Post("/check", func(rw http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				err = r.Body.Close()
				require.NoError(t, err)

				assert.Equal(t, []byte(test.body), body)

				_, err = rw.Write(body)
				require.NoError(t, err)

				rw.WriteHeader(http.StatusOK)
			})

			testServer := httptest.NewServer(r)
			defer testServer.Close()

			req, err := http.NewRequest(http.MethodPost, testServer.URL+"/check", bytes.NewBuffer([]byte(test.body)))
			assert.NoError(t, err)

			req.Header.Set("X-Real-IP", test.xRealIP)

			r.ServeHTTP(respRecorder, req)
			resp := respRecorder.Result()

			require.Equal(t, test.expectedStatus, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, test.expectedBody, body)

			err = resp.Body.Close()
			assert.NoError(t, err)
		})
	}
}
