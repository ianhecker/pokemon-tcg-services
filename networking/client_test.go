package networking_test

import (
	"context"
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

func newLazyTestServer(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/sleeping", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5)
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

			ctx := context.Background()
			logger := zap.NewNop().Sugar()

			client := networking.NewClient(logger)
			body, status, err := client.Get(ctx, srv.URL+"/hello")

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

	t.Run("cancels with context", func(t *testing.T) {
		srv := newLazyTestServer(t)
		defer srv.Close()

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		logger := zap.NewNop().Sugar()

		client := networking.NewClient(logger)

		errorChan := make(chan error, 1)
		go func() {
			_, _, err := client.Get(ctx, srv.URL+"/sleeping")
			errorChan <- err
		}()

		cancel()

		select {
		case err := <-errorChan:
			assert.ErrorContains(t, err, "context canceled")

		case <-time.After(1 * time.Second):
			assert.Fail(t, "test timed out after 1 sec")
		}

	})
}
