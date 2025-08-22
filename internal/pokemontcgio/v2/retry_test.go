package v2_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v2 "github.com/ianhecker/pokemon-tcg-services/internal/pokemontcgio/v2"
	"github.com/ianhecker/pokemon-tcg-services/internal/retry"
)

func TestV2_RetryForStatus(t *testing.T) {
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
			state := v2.RetryForStatus(test.status)
			assert.Equal(t, test.state, state,
				fmt.Sprintf("status:%d want:%s got:%s", test.status, test.state, state),
			)
		})
	}
}

func wrap(e error) error { return fmt.Errorf("wrap: %w", e) }

func TestV2_RetryForError(t *testing.T) {
	tests := []struct {
		name  string
		error error
		state retry.RetryState
	}{
		{"canceled", context.Canceled, retry.Fail},
		{"wrapped canceled", wrap(context.Canceled), retry.Fail},
		{"joined canceled", errors.Join(context.Canceled, io.EOF), retry.Fail},

		{"deadline", context.DeadlineExceeded, retry.Backoff},
		{"wrapped deadline", wrap(context.DeadlineExceeded), retry.Backoff},

		{"net timeout", &net.DNSError{IsTimeout: true}, retry.Backoff},
		{"net op timeout", &net.OpError{Op: "dial", Net: "tcp", Err: &net.DNSError{IsTimeout: true}}, retry.Backoff},

		{"other", errors.New("other"), retry.Backoff},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			state := v2.RetryForError(test.error)
			require.Equal(t, test.state, state)
		})
	}
}
