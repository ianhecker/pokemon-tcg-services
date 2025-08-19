package pokemontcgio_test

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
		{"408 retry no backoff", 408, retry.RetryNoBackoff},
		{"425 retry no backoff", 425, retry.RetryNoBackoff},
		// 5xx
		{"500 retry no backoff", 500, retry.RetryNoBackoff},
		{"502 retry no backoff", 502, retry.RetryNoBackoff},
		{"504 retry no backoff", 504, retry.RetryNoBackoff},
		// 5xx - catch all
		{"501 retry no backoff", 501, retry.RetryNoBackoff},
		{"599 retry no backoff", 599, retry.RetryNoBackoff},

		// retry - with backoff
		// 4xx
		{"429 retry w/ backoff", 429, retry.RetryWithBackoff},
		// 5xx
		{"503 retry w/ backoff", 503, retry.RetryWithBackoff},

		// no retry
		// 2xx
		{"2xx no retry", 200, retry.NoRetry},
		// 3xx
		{"3xx no retry", 300, retry.NoRetry},
		// 4xx
		{"4xx no retry", 400, retry.NoRetry},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := v2.RetryForStatus(test.status)
			assert.Equal(t, test.state, state)
		})
	}
}
