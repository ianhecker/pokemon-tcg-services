package pokemontcgio

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/ianhecker/pokemon-tcg-services/networking"
	v2 "github.com/ianhecker/pokemon-tcg-services/networking/pokemontcgio/v2"
	"github.com/ianhecker/pokemon-tcg-services/retry"
)

type Client struct {
	client  networking.ClientInterface
	timeout time.Duration
}

func NewClient(
	client networking.ClientInterface,
	timeout time.Duration,
) *Client {
	return &Client{
		client:  client,
		timeout: timeout,
	}
}

type Result struct {
	Body   []byte
	Status int
	Err    error
}

func (c *Client) MakeRetryFunc(url string) (*Result, retry.RetryFunc) {

	result := &Result{}

	retryFunc := func(ctx context.Context) (retry.RetryState, error) {
		attemptCtx, cancel := context.WithTimeout(ctx, c.timeout)
		defer cancel()

		r := result
		r.Body, r.Status, r.Err = c.client.Get(attemptCtx, url)

		err := r.Err
		status := r.Status

		if errors.Is(err, context.DeadlineExceeded) {
			return retry.RetryWithBackoff, nil
		}

		if errors.Is(err, context.Canceled) {
			return retry.NoRetry, err
		}

		var netError net.Error
		if errors.As(err, &netError) && netError.Timeout() {
			return retry.RetryWithBackoff, nil
		}

		if status != 0 {
			switch v2.RetryForStatus(status) {
			case retry.NoRetry:
				return retry.NoRetry, err
			case retry.RetryNoBackoff:
				return retry.RetryNoBackoff, nil
			case retry.RetryWithBackoff:
				return retry.RetryWithBackoff, nil
			}
		}
		return retry.NoRetry, err
	}
	return result, retryFunc
}
