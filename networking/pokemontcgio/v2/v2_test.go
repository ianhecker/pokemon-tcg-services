package v2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	v2 "github.com/ianhecker/pokemon-tcg-services/networking/pokemontcgio/v2"
	"github.com/ianhecker/pokemon-tcg-services/retry"
)

func TestV2_RetryForStatus(t *testing.T) {
	tests := []struct {
		name   string
		status int
		state  retry.RetryState
	}{
		// retry - no backoff
		// 4xx
		{"408 retry no backoff", 408, retry.Yes},
		{"425 retry no backoff", 425, retry.Yes},
		// 5xx
		{"500 retry no backoff", 500, retry.Yes},
		{"502 retry no backoff", 502, retry.Yes},
		{"504 retry no backoff", 504, retry.Yes},
		// 5xx - catch all
		{"501 retry no backoff", 501, retry.Yes},
		{"599 retry no backoff", 599, retry.Yes},

		// retry - with backoff
		// 4xx
		{"429 retry w/ backoff", 429, retry.WithBackoff},
		// 5xx
		{"503 retry w/ backoff", 503, retry.WithBackoff},

		// no retry
		// 2xx
		{"2xx no retry", 200, retry.No},
		// 3xx
		{"3xx no retry", 300, retry.No},
		// 4xx
		{"4xx no retry", 400, retry.No},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := v2.RetryForStatus(test.status)
			assert.Equal(t, test.state, state)
		})
	}
}
