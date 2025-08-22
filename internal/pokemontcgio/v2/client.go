package v2

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/ianhecker/pokemon-tcg-services/internal/networking"
	"github.com/ianhecker/pokemon-tcg-services/internal/retry"
)

type APIClientInterface interface {
	MakeRetryFunc(url string) (*Result, retry.RetryFunc)
	MakeRetryable(url string) (*Result, retry.RetryableInterface)
}

type Client struct {
	log     *zap.SugaredLogger
	client  networking.HttpClientInterface
	timeout time.Duration
}

func NewClient(logger *zap.SugaredLogger) APIClientInterface {
	httpClient := networking.NewClient(logger)

	return &Client{
		log:     logger,
		client:  httpClient,
		timeout: AttemptTimeout,
	}
}

func NewClientFromRaw(
	logger *zap.SugaredLogger,
	client networking.HttpClientInterface,
	timeout time.Duration,
) APIClientInterface {
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

func (c *Client) MakeRetryable(url string) (*Result, retry.RetryableInterface) {

	result, retryFunc := c.MakeRetryFunc(url)
	retryable := retry.MakeRetryable(Retries, BackoffInSeconds, retryFunc)
	return result, retryable
}
