package v2

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/ianhecker/pokemon-tcg-services/networking"
	"github.com/ianhecker/pokemon-tcg-services/retry"
	"go.uber.org/zap"
)

type Client struct {
	log     *zap.SugaredLogger
	client  networking.ClientInterface
	timeout time.Duration
}

func NewClient(
	logger *zap.SugaredLogger,
	client networking.ClientInterface,
	timeout time.Duration,
) *Client {
	return &Client{
		log:     logger,
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
			c.log.Infow("retry deadline met", "url", url, "backoff", c.timeout, "err", err)
			return retry.WithBackoff, nil
		}

		if errors.Is(err, context.Canceled) {
			c.log.Errorw("retry canceled", "url", url, "err", err)
			return retry.No, err
		}

		var netError net.Error
		if errors.As(err, &netError) && netError.Timeout() {
			c.log.Errorw("retry network error", "url", url, "backoff", c.timeout, "err", err)
			return retry.WithBackoff, nil
		}

		if status != 0 {
			switch RetryForStatus(status) {
			case retry.No:
				c.log.Infow("no retry for status", "url", url, "status", status, "err", err)
				return retry.No, err
			case retry.Yes:
				c.log.Infow("retry for status", "url", url, "status", status, "err", err)
				return retry.Yes, nil
			case retry.WithBackoff:
				c.log.Infow("retry with backoff for status", "url", url, "status", status, "err", err)
				return retry.WithBackoff, nil
			}
		}
		return retry.No, err
	}
	return result, retryFunc
}
