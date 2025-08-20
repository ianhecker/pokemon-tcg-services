package v2

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/ianhecker/pokemon-tcg-services/networking"
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
			return retry.WithBackoff, nil
		}

		if errors.Is(err, context.Canceled) {
			return retry.No, err
		}

		var netError net.Error
		if errors.As(err, &netError) && netError.Timeout() {
			return retry.WithBackoff, nil
		}

		if status != 0 {
			switch RetryForStatus(status) {
			case retry.No:
				return retry.No, err
			case retry.Yes:
				return retry.Yes, nil
			case retry.WithBackoff:
				return retry.WithBackoff, nil
			}
		}
		return retry.No, err
	}
	return result, retryFunc
}
