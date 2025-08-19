package networking_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/ianhecker/pokemon-tcg-services/networking"
)

func newTestServer(t *testing.T, status int, body string) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if body != "" {
			fmt.Fprint(w, body)
		}
	})

	return httptest.NewServer(mux)
}

func TestClient_Get(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
	}{
		// Success cases
		{"OK success", 200, `{"data":{"id":"xy7-54","name":"Gardevoir"}}`},

		// Client errors (4xx)
		{"Bad Request", 400, `{"error":"invalid ID format"}`},
		{"Unauthorized", 401, `{"error":"missing API key"}`},
		{"Forbidden", 403, `{"error":"access denied"}`},
		{"Not Found", 404, `{"error":"card not found"}`},
		{"Unprocessable Entity", 422, `{"error":"validation failed"}`},

		// // Rate limiting / throttling
		{"Too Many Requests", 429, `{"error":"rate limited","retry_after":30}`},

		// // Server errors (5xx)
		{"Internal Server Error", 500, `{"error":"server crashed"}`},
		{"Bad Gateway", 502, `{"error":"upstream failed"}`},
		{"Service Unavailable", 503, `{"error":"overloaded"}`},
		{"Gateway Timeout", 504, `{"error":"timeout"}`},

		// // Catch-all unexpected code
		{"Weird code", 599, `{"error":"unknown"}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srv := newTestServer(t, test.status, test.body)
			defer srv.Close()

			logger := zap.NewNop().Sugar()
			timeout := 3 * time.Second
			client := networking.NewClient(logger, timeout)
			body, status, err := client.Get(srv.URL + "/hello")

			assert.Equal(t, test.status, status)
			assert.Equal(t, string(test.body), string(body))

			if 200 <= test.status && test.status < 400 {
				assert.NoError(t, err)

			} else {
				expected := fmt.Sprintf("unexpected status code: %d", test.status)
				assert.ErrorContains(t, err, expected)
			}
		})
	}
}
