package server

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	tests := []struct {
		name                string
		fnc                 func(code int, body string) http.HandlerFunc
		addr                string
		path                string
		body                string
		expectedBody        string
		expectedStatusCode  int
		expectedShutdownErr error
	}{
		{
			name: "Check run server and process request",
			fnc: func(code int, body string) http.HandlerFunc {
				return func(rw http.ResponseWriter, r *http.Request) {
					rw.WriteHeader(code)
					rw.Write([]byte(body))
				}
			},
			addr:               "localhost:8082",
			path:               "/test/server/",
			body:               "my test body",
			expectedBody:       "my test body",
			expectedStatusCode: 200,
		},
	}

	testClient := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.Handle(test.path, http.HandlerFunc(test.fnc(test.expectedStatusCode, test.expectedBody)))

			srv := New(context.Background(), test.addr, mux)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := srv.Run()
				if errors.Is(err, http.ErrServerClosed) {
					return
				}
				assert.NoError(t, err)
			}()

			body := bytes.NewBuffer([]byte(test.body))
			req, err := http.NewRequest(http.MethodPost, "http://"+test.addr+test.path, body)
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

			err = srv.Shutdown(context.Background())
			assert.NoError(t, err)

			wg.Wait()
		})
	}
}
