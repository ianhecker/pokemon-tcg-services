package retry

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type RetryableInterface interface {
	Retry(ctx context.Context) (bool, error)
	RetriesRemaining() int
}

type RetryState int

const (
	No RetryState = iota
	Yes
	WithBackoff
)

type RetryFunc func(ctx context.Context) (RetryState, error)

type Retryable struct {
	Attempt int
	Retries int
	State   RetryState
	Sleeper Sleepable
	Do      RetryFunc
}

func MakeRetryable(
	retries int,
	duration time.Duration,
	do RetryFunc,
) RetryableInterface {
	return &Retryable{
		Retries: retries,
		Sleeper: NewSleeper(duration),
		Do:      do,
		State:   Yes,
	}
}

func MakeRetryableFromRaw(
	attempt int,
	retries int,
	state RetryState,
	sleeper Sleepable,
	do RetryFunc,
) RetryableInterface {
	return &Retryable{
		Attempt: attempt,
		Retries: retries,
		State:   state,
		Sleeper: sleeper,
		Do:      do,
	}
}

func (r *Retryable) Retry(ctx context.Context) (bool, error) {
	if r.State == No {
		return false, errors.New("starting state was not retryable")
	}
	if r.Retries <= 0 {
		return false, errors.New("retries were zero")
	}
	if r.Do == nil {
		return false, errors.New("nil Do function")
	}

	r.Retries--
	r.Attempt++

	if r.State == WithBackoff && r.Attempt > 1 && r.Sleeper.Duration() > 0 {
		stop, done := r.Sleeper.Sleep()
		defer stop()

		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-done:
		}
	}

	newState, err := r.Do(ctx)
	r.State = newState

	if newState == No || err != nil {
		return false, err
	}
	return true, nil
}

func (retry Retryable) RetriesRemaining() int {
	return retry.Retries
}

func RunRetryable(ctx context.Context, retryable RetryableInterface) error {
	retries := retryable.RetriesRemaining()

	for i := retries; i > 0; i-- {

		retry, err := retryable.Retry(ctx)
		if err != nil {
			return err

		} else if retry != true {
			return nil
		}
	}
	return fmt.Errorf("retried: %d times. Out of retries", retries)
}
