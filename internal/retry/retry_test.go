package retry_test

import (
	"testing"

	"github.com/ianhecker/pokemon-tcg-services/internal/retry"
	"github.com/stretchr/testify/assert"
)

func TestRetryState_String(t *testing.T) {
	var tests = []struct {
		name     string
		state    retry.RetryState
		expected string
	}{
		{"fail", retry.Fail, "fail"},
		{"complete", retry.Complete, "complete"},
		{"yes", retry.Yes, "yes"},
		{"backoff", retry.Backoff, "backoff"},
		{"unknown", 1234, "unknown state: 1234"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.state.String()
			assert.Equal(t, test.expected, got)
		})
	}
}
