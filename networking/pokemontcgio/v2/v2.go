package pokemontcgio

import (
	"fmt"
	"net/http"

	"github.com/ianhecker/pokemon-tcg-services/retry"
)

const CardByIDFmt string = "api.pokemontcg.io/v2/cards/%s"

func CardByID(ID string) string {
	return fmt.Sprintf(CardByIDFmt, ID)
}

func RetryForStatus(status int) retry.RetryState {
	switch status {

	case http.StatusRequestTimeout, // 408
		http.StatusTooEarly,            // 425
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusGatewayTimeout:      // 504

		return retry.RetryNoBackoff

	case http.StatusTooManyRequests, // 429
		http.StatusServiceUnavailable: // 503

		return retry.RetryWithBackoff
	}

	if 500 <= status && status <= 599 {
		return retry.RetryNoBackoff
	}
	return retry.NoRetry
}
