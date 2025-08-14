package networking_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ianhecker/pokemon-tcg-services/networking"
	"github.com/stretchr/testify/assert"
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

func TestClient_NewRequest(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		status := 200
		body := `{"message":"hi"}`
		srv := newTestServer(t, status, body)
		defer srv.Close()

		client := networking.NewClient()
		got, err := client.Get(srv.URL + "/hello")
		assert.NoError(t, err)
		assert.Equal(t, string(body), string(got))
	})

	t.Run("status code is not OK", func(t *testing.T) {
		status := 400
		body := ``
		srv := newTestServer(t, status, body)
		defer srv.Close()

		client := networking.NewClient()
		_, err := client.Get(srv.URL + "/hello")
		assert.Error(t, err, fmt.Sprintf("unexpected status code: %d", status))
	})
}
