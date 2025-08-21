package v2

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/ianhecker/pokemon-tcg-services/retry"
)

func RetryForStatus(status int) retry.RetryState {
	if 200 <= status && status <= 299 {
		return retry.Complete
	}

	switch status {
	case http.StatusNotModified: // 304
		return retry.Complete

	case http.StatusRequestTimeout, // 408
		http.StatusTooEarly,            // 425
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusGatewayTimeout:      // 504
		return retry.Yes

	case http.StatusTooManyRequests, // 429
		http.StatusServiceUnavailable: // 503
		return retry.Backoff
	}

	if 500 <= status && status <= 599 {
		return retry.Yes
	}
	return retry.Fail
}

func RetryForError(err error) retry.RetryState {
	if errors.Is(err, context.Canceled) {
		return retry.Fail
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return retry.Backoff
	}

	var netError net.Error
	if errors.As(err, &netError) && netError.Timeout() {
		return retry.Backoff
	}
	return retry.Backoff
}
