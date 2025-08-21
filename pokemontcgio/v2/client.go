package v2

import (
	"context"
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

		if err != nil {
			state := RetryForError(err)
			switch state {
			case retry.Backoff:
				c.log.Infow("retry with backoff for error", "url", url, "status", status, "err", err)
				return state, nil
			case retry.Fail:
				c.log.Infow("retry fail for error", "url", url, "status", status, "err", err)
				return state, err
			}
		}

		if status != 0 {
			state := RetryForStatus(status)
			switch state {
			case retry.Complete:
				c.log.Infow("retry completed", "url", url, "status", status, "err", err)
				return state, nil
			case retry.Yes:
				c.log.Infow("retry for status", "url", url, "status", status, "err", err)
				return state, err
			case retry.Backoff:
				c.log.Infow("retry with backoff for status", "url", url, "status", status, "err", err)
				return state, err
			case retry.Fail:
				c.log.Infow("no retry for status", "url", url, "status", status, "err", err)
				return state, err
			}
		}
		return retry.Fail, err
	}
	return result, retryFunc
}
