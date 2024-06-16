package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareCompress(t *testing.T) {
	tests := []struct {
		name                    string
		acceptEncoding          string
		path                    string
		body                    io.Reader
		expectedContentEncoding string
	}{
		{
			name:                    "check compress application/json",
			acceptEncoding:          "gzip",
			path:                    "/get_application_json",
			expectedContentEncoding: "gzip",
		},
		{
			name:                    "check compress text/html",
			acceptEncoding:          "gzip",
			path:                    "/get_text_html",
			expectedContentEncoding: "gzip",
		},
		{
			name:                    "check no compress for no Accept-Encoding",
			acceptEncoding:          "",
			path:                    "/get_application_json",
			expectedContentEncoding: "",
		},
		{
			name:                    "check no compress for unknown Content-Type",
			acceptEncoding:          "gzip",
			path:                    "/get_unknown_content_type",
			expectedContentEncoding: "",
		},
	}

	r := chi.NewRouter()
	r.Use(Compress)

	r.Get("/get_application_json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte("textstring"))
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/get_text_html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, err := w.Write([]byte("textstring"))
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/get_unknown_content_type", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "unknown")
		_, err := w.Write([]byte("textstring"))
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	tc := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, ts.URL+test.path, nil)
			assert.NoError(t, err)

			req.Header.Set("Accept-Encoding", test.acceptEncoding)
			assert.NoError(t, err)

			resp, err := tc.Do(req)
			assert.NoError(t, err)
			//nolint:errcheck // not need check error in test
			defer resp.Body.Close()

			var reader io.ReadCloser

			switch resp.Header.Get("Content-Encoding") {
			case "gzip":
				var err error
				reader, err = gzip.NewReader(resp.Body)
				assert.NoError(t, err)
			default:
				reader = resp.Body
			}

			respBody, err := io.ReadAll(reader)
			assert.NoError(t, err)

			err = reader.Close()
			assert.NoError(t, err)

			assert.Equal(t, "textstring", string(respBody),
				"response text doesn't match; expected:%q, got:%q", "textstring", string(respBody))

			respContentEncoding := resp.Header.Get("Content-Encoding")
			assert.Equal(t, test.expectedContentEncoding, respContentEncoding,
				"expected Content-Encoding %q but got %q", test.expectedContentEncoding, respContentEncoding)
		})
	}
}
