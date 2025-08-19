package retry

import (
	"context"
	"errors"
	"time"
)

type RetryableInterface interface {
	Retry(ctx context.Context) (bool, error)
}

type RetryState int

const (
	NoRetry RetryState = iota
	RetryNoBackoff
	RetryWithBackoff
)

type RetryFunc func(ctx context.Context) (RetryState, error)

type Retryable struct {
	Attempt int
	Retries int
	State   RetryState
	Sleep   time.Duration
	Do      RetryFunc
	Body    []byte
}

func MakeRetryable(
	retries int,
	sleep time.Duration,
	do RetryFunc,
) RetryableInterface {
	return &Retryable{
		Retries: retries,
		Sleep:   sleep,
		Do:      do,
		State:   RetryNoBackoff,
	}
}

func (r *Retryable) Retry(ctx context.Context) (bool, error) {
	if r.State == NoRetry {
		return false, errors.New("retry state was not retriable")
	}
	if r.Retries <= 0 {
		return false, errors.New("retries were zero")
	}
	if r.Do == nil {
		return false, errors.New("nil Do function")
	}

	r.Retries--
	r.Attempt++

	if r.State == RetryWithBackoff && r.Attempt > 1 && r.Sleep > 0 {
		timer := time.NewTimer(r.Sleep)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-timer.C:
		}
	}

	newState, err := r.Do(ctx)
	r.State = newState

	if newState == NoRetry || err != nil {
		return false, err
	}
	return true, nil
}

func Do(ctx context.Context, retryable RetryableInterface) error {
	for {
		retry, err := retryable.Retry(ctx)
		if err != nil {
			return err

		} else if retry != true {
			return nil
		}
	}
}
