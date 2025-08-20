package v2

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/ianhecker/pokemon-tcg-services/retry"
)

func RetryForStatus(status int) retry.RetryState {
	switch status {

	case http.StatusRequestTimeout, // 408
		http.StatusTooEarly,            // 425
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusGatewayTimeout:      // 504

		return retry.Yes

	case http.StatusTooManyRequests, // 429
		http.StatusServiceUnavailable: // 503

		return retry.WithBackoff
	}

	if 500 <= status && status <= 599 {
		return retry.Yes
	}
	return retry.No
}

func RetryForError(err error) retry.RetryState {
	if err == nil {
		return retry.No
	}

	if errors.Is(err, context.Canceled) {
		return retry.No
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return retry.WithBackoff
	}

	var netError net.Error
	if errors.As(err, &netError) && netError.Timeout() {
		return retry.WithBackoff
	}
	return retry.WithBackoff
}
