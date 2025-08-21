package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ianhecker/pokemon-tcg-services/internal/mocks"
	"github.com/ianhecker/pokemon-tcg-services/retry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRetryable_Retry(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		retries := 1
		duration := time.Duration(0)
		calls := 0
		do := func(ctx context.Context) (retry.RetryState, error) {
			calls++
			return retry.No, nil
		}

		ctx := context.Background()
		retry := retry.MakeRetryable(retries, duration, do)

		repeat, err := retry.Retry(ctx)
		assert.NoError(t, err)
		assert.Equal(t, false, repeat)
		assert.Equal(t, 1, calls)
	})

	t.Run("happy path retry", func(t *testing.T) {
		retries := 2
		duration := time.Duration(0)
		calls := 0
		do := func(ctx context.Context) (retry.RetryState, error) {
			calls++
			if calls == 1 {
				return retry.Yes, nil
			} else {
				return retry.No, nil
			}
		}

		ctx := context.Background()
		retry := retry.MakeRetryable(retries, duration, do)

		repeat, err := retry.Retry(ctx)
		assert.NoError(t, err)
		assert.Equal(t, true, repeat)
		assert.Equal(t, 1, calls)

		repeat, err = retry.Retry(ctx)
		assert.NoError(t, err)
		assert.Equal(t, false, repeat)
		assert.Equal(t, 2, calls)
	})

	t.Run("happy path backoff", func(t *testing.T) {
		retries := 1
		calls := 0
		do := func(ctx context.Context) (retry.RetryState, error) {
			calls++
			return retry.No, nil
		}

		stop := func() bool { return true }
		done := make(chan time.Time, 1)
		sleepableMock := mocks.NewMockSleepable(t)
		sleepableMock.
			On("Duration").
			Return(1 * time.Millisecond).
			Once()
		sleepableMock.
			On("Sleep").
			Return(stop, (<-chan time.Time)(done)).
			Once().
			Run(func(args mock.Arguments) {
				time.Sleep(time.Millisecond)
				close(done)
			})

		ctx := context.Background()
		retry := retry.MakeRetryableFromRaw(
			1,
			retries,
			retry.WithBackoff,
			sleepableMock,
			do,
		)

		repeat, err := retry.Retry(ctx)
		assert.NoError(t, err)
		assert.Equal(t, false, repeat)
		assert.Equal(t, 1, calls)
		sleepableMock.AssertExpectations(t)
	})

	t.Run("context canceled in backoff", func(t *testing.T) {
		retries := 1
		calls := 0
		do := func(ctx context.Context) (retry.RetryState, error) {
			calls++
			return retry.No, nil
		}

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)

		stop := func() bool { return true }
		done := make(chan time.Time, 1)
		sleepableMock := mocks.NewMockSleepable(t)
		sleepableMock.
			On("Duration").
			Return(1 * time.Millisecond).
			Once()
		sleepableMock.
			On("Sleep").
			Return(stop, (<-chan time.Time)(done)).
			Once().
			Run(func(args mock.Arguments) {
				time.Sleep(time.Millisecond)
				cancel()
			})

		retry := retry.MakeRetryableFromRaw(
			1,
			retries,
			retry.WithBackoff,
			sleepableMock,
			do,
		)

		repeat, err := retry.Retry(ctx)
		assert.ErrorContains(t, err, "context canceled")
		assert.Equal(t, false, repeat)
		assert.Equal(t, 0, calls)
		sleepableMock.AssertExpectations(t)
	})

	t.Run("retry returns error", func(t *testing.T) {
		retries := 1
		duration := time.Duration(0)
		calls := 0
		expected := errors.New("error")
		do := func(ctx context.Context) (retry.RetryState, error) {
			calls++
			return retry.No, expected
		}

		ctx := context.Background()
		retry := retry.MakeRetryable(retries, duration, do)

		repeat, err := retry.Retry(ctx)
		assert.ErrorIs(t, err, expected)
		assert.Equal(t, false, repeat)
		assert.Equal(t, 1, calls)
	})

	t.Run("state starts at no retry", func(t *testing.T) {
		ctx := context.Background()
		retry := retry.MakeRetryableFromRaw(
			0,
			0,
			retry.No,
			nil,
			nil,
		)

		repeat, err := retry.Retry(ctx)
		assert.EqualError(t, err, "starting state was not retryable")
		assert.Equal(t, false, repeat)
	})

	t.Run("given nil function", func(t *testing.T) {
		ctx := context.Background()
		retry := retry.MakeRetryableFromRaw(
			0,
			1,
			retry.Yes,
			nil,
			nil,
		)

		repeat, err := retry.Retry(ctx)
		assert.EqualError(t, err, "nil Do function")
		assert.Equal(t, false, repeat)
	})

	t.Run("retries are zero", func(t *testing.T) {
		ctx := context.Background()
		retry := retry.MakeRetryableFromRaw(
			0,
			0,
			retry.Yes,
			nil,
			nil,
		)

		repeat, err := retry.Retry(ctx)
		assert.EqualError(t, err, "retries were zero")
		assert.Equal(t, false, repeat)
	})
}

func TestRunRetryable(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		retryable := mocks.NewMockRetryableInterface(t)
		retryable.
			On("RetriesRemaining").
			Return(1).
			Once()
		retryable.
			On("Retry", ctx).
			Return(false, nil).
			Once()

		err := retry.RunRetryable(ctx, retryable)
		assert.NoError(t, err)
		retryable.AssertExpectations(t)
	})

	t.Run("happy path retry once", func(t *testing.T) {
		ctx := context.Background()

		retryable := mocks.NewMockRetryableInterface(t)
		retryable.
			On("RetriesRemaining").
			Return(2).
			Once()
		retryable.
			On("Retry", ctx).
			Return(true, nil).
			Once()
		retryable.
			On("Retry", ctx).
			Return(false, nil).
			Once()

		err := retry.RunRetryable(ctx, retryable)
		assert.NoError(t, err)
		retryable.AssertExpectations(t)
	})

	t.Run("retry error", func(t *testing.T) {
		ctx := context.Background()
		expected := errors.New("error")

		retryable := mocks.NewMockRetryableInterface(t)
		retryable.
			On("RetriesRemaining").
			Return(1).
			Once()
		retryable.
			On("Retry", ctx).
			Return(false, expected)

		err := retry.RunRetryable(ctx, retryable)
		assert.ErrorIs(t, err, expected)
	})

	t.Run("retry no op", func(t *testing.T) {
		ctx := context.Background()

		retryable := mocks.NewMockRetryableInterface(t)
		retryable.
			On("RetriesRemaining").
			Return(0).
			Once()

		err := retry.RunRetryable(ctx, retryable)
		assert.ErrorContains(t, err, "retried: 0 times. Out of retries")
		retryable.AssertExpectations(t)
	})

	t.Run("ran out of retries", func(t *testing.T) {
		ctx := context.Background()
		retries := 3

		retryable := mocks.NewMockRetryableInterface(t)
		retryable.
			On("RetriesRemaining").
			Return(retries).
			Once()
		retryable.
			On("Retry", ctx).
			Return(true, nil).
			Times(retries)

		err := retry.RunRetryable(ctx, retryable)
		assert.ErrorContains(t, err, "retried: 3 times. Out of retries")
	})
}
