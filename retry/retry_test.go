package retry_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ianhecker/pokemon-tcg-services/internal/mocks"
	"github.com/ianhecker/pokemon-tcg-services/retry"
	"github.com/stretchr/testify/assert"
)

func TestRetryable_Retry(t *testing.T) {
	t.Run("", func(t *testing.T) {})
	t.Run("", func(t *testing.T) {})
	t.Run("", func(t *testing.T) {})
}

func TestDo(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		retryable := mocks.NewMockRetryableInterface(t)
		retryable.On("Retry", ctx).Return(false, nil)

		err := retry.Do(ctx, retryable)
		assert.NoError(t, err)
	})

	t.Run("retry", func(t *testing.T) {
		ctx := context.Background()

		retryable := mocks.NewMockRetryableInterface(t)
		retryable.On("Retry", ctx).Return(true, nil)

		err := retry.Do(ctx, retryable)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		ctx := context.Background()

		expected := errors.New("error")

		retryable := mocks.NewMockRetryableInterface(t)
		retryable.On("Retry", ctx).Return(false, expected)

		err := retry.Do(ctx, retryable)
		assert.ErrorIs(t, err, expected)
	})
}
