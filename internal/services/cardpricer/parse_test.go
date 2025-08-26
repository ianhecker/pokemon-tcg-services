package cardpricer_test

import (
	"net/http"
	"testing"

	"github.com/ianhecker/pokemon-tcg-services/internal/services/cardpricer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRequest(t *testing.T, query string) *http.Request {
	request, err := http.NewRequest(http.MethodGet, "localhost:8080/v1/cards"+query, nil)
	require.NoError(t, err)
	return request
}

func TestParseCardQuery(t *testing.T) {
	var tests = []struct {
		name     string
		request  *http.Request
		expected string
		err      string
	}{
		{"happy path", newRequest(t, "?id=base1-1"), "base1-1", ""},
		{"request is nil", nil, "", "request is nil"},
		{"missing query", newRequest(t, ""), "", "missing required query: id"},
		{"empty query", newRequest(t, "?id="), "", "invalid card ID: ''"},
		{"bad ID", newRequest(t, "?id=bad ID"), "", "invalid card ID: 'bad ID'"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ID, err := cardpricer.ParseCardQuery(test.request)

			if test.err == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, ID.String())
			} else {
				assert.ErrorContains(t, err, test.err)
			}
		})
	}
}
