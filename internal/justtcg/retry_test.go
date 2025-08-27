package justtcg_test

import (
	"fmt"
	"testing"

	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg"
	"github.com/ianhecker/pokemon-tcg-services/internal/retry"
	"github.com/stretchr/testify/assert"
)

func TestRetryForStatus(t *testing.T) {
	tests := []struct {
		name   string
		status int
		state  retry.RetryState
	}{
		// retry - complete
		{"2xx complete", 200, retry.Complete},
		{"304 complete", 304, retry.Complete},
		// retry - yes
		{"408 retry no backoff", 408, retry.Yes},
		{"425 retry no backoff", 425, retry.Yes},
		{"500 retry no backoff", 500, retry.Yes},
		{"502 retry no backoff", 502, retry.Yes},
		{"504 retry no backoff", 504, retry.Yes},
		// retry - backoff
		{"429 retry w/ backoff", 429, retry.Backoff},
		{"503 retry w/ backoff", 503, retry.Backoff},
		// retry - yes - catchall
		{"501 retry no backoff", 501, retry.Yes},
		{"599 retry no backoff", 599, retry.Yes},
		// retry -fail
		{"300 complete", 300, retry.Fail},
		{"no retry", 0, retry.Fail},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := justtcg.RetryForStatus(test.status)
			assert.Equal(t, test.state, state,
				fmt.Sprintf("status:%d want:%s got:%s", test.status, test.state, state),
			)
		})
	}
}
