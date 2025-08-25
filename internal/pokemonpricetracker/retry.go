package pokemonpricetracker

import (
	"net/http"

	"github.com/ianhecker/pokemon-tcg-services/internal/retry"
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
