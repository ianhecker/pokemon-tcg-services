package v2_test

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/ianhecker/pokemon-tcg-services/internal/mocks"
	"github.com/ianhecker/pokemon-tcg-services/internal/retry"
	v2 "github.com/ianhecker/pokemon-tcg-services/pokemontcgio/v2"
)

func TestClient_MakeRetryFunc(t *testing.T) {
	t.Run("deadline exceeded", func(t *testing.T) {
		url := "url"
		mockClient := mocks.NewMockAPIClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(nil, 0, context.DeadlineExceeded)

		logger := zap.NewNop().Sugar()
		timeout := 100 * time.Millisecond

		client := v2.NewClientFromRaw(logger, mockClient, timeout)
		result, retryFunc := client.MakeRetryFunc(url)

		ctx, cancel := context.WithTimeout(context.Background(), 1)
		defer cancel()

		state, err := retryFunc(ctx)
		assert.NoError(t, err)
		assert.Equal(t, retry.Backoff, state)
		assert.ErrorIs(t, result.Err, context.DeadlineExceeded)
	})

	t.Run("canceled", func(t *testing.T) {
		url := "url"

		ctx, cancel := context.WithTimeout(context.Background(), 1)

		mockClient := mocks.NewMockAPIClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(nil, 0, context.Canceled).
			Run(func(args mock.Arguments) {
				time.Sleep(1 * time.Millisecond)
				cancel()
			})

		logger := zap.NewNop().Sugar()
		timeout := 100 * time.Millisecond

		client := v2.NewClientFromRaw(logger, mockClient, timeout)
		result, retryFunc := client.MakeRetryFunc(url)

		state, err := retryFunc(ctx)
		assert.ErrorIs(t, err, context.Canceled)
		assert.Equal(t, retry.Fail, state)
		assert.ErrorIs(t, result.Err, context.Canceled)
	})

	t.Run("networking error", func(t *testing.T) {
		url := "url"

		var netError error = &net.OpError{
			Op:  "dial",
			Net: "tcp",
			Err: &net.DNSError{Err: "i/o timeout", IsTimeout: true},
		}

		mockClient := mocks.NewMockAPIClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(nil, 0, netError)

		logger := zap.NewNop().Sugar()
		timeout := 100 * time.Millisecond

		client := v2.NewClientFromRaw(logger, mockClient, timeout)
		result, retryFunc := client.MakeRetryFunc(url)

		ctx, cancel := context.WithTimeout(context.Background(), 1)
		defer cancel()

		state, err := retryFunc(ctx)
		assert.NoError(t, err)
		assert.Equal(t, retry.Backoff, state)
		assert.ErrorAs(t, result.Err, &netError)
	})

	t.Run("other error", func(t *testing.T) {
		url := "url"
		otherErr := errors.New("other")

		mockClient := mocks.NewMockAPIClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(nil, 0, otherErr)

		logger := zap.NewNop().Sugar()
		timeout := 100 * time.Millisecond

		client := v2.NewClientFromRaw(logger, mockClient, timeout)
		result, retryFunc := client.MakeRetryFunc(url)

		ctx, cancel := context.WithTimeout(context.Background(), 1)
		defer cancel()

		state, err := retryFunc(ctx)
		assert.NoError(t, err)
		assert.Equal(t, retry.Backoff, state)
		assert.ErrorIs(t, result.Err, otherErr)
	})

	t.Run("status no retry", func(t *testing.T) {
		body := []byte(`{"message":"hi"}`)
		status := 200
		url := "url"
		mockClient := mocks.NewMockAPIClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(body, status, nil)

		logger := zap.NewNop().Sugar()
		timeout := 100 * time.Millisecond

		client := v2.NewClientFromRaw(logger, mockClient, timeout)
		result, retryFunc := client.MakeRetryFunc(url)

		state, err := retryFunc(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, retry.Complete, state)

		assert.Equal(t, string(body), string(result.Body))
		assert.Equal(t, status, result.Status)
		assert.NoError(t, result.Err)
	})

	t.Run("status retry", func(t *testing.T) {
		body := []byte(`{"message":"hi"}`)
		status := 504
		url := "url"
		mockClient := mocks.NewMockAPIClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(body, status, nil)

		logger := zap.NewNop().Sugar()
		timeout := 100 * time.Millisecond

		client := v2.NewClientFromRaw(logger, mockClient, timeout)
		result, retryFunc := client.MakeRetryFunc(url)

		state, err := retryFunc(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, retry.Yes, state)

		assert.Equal(t, string(body), string(result.Body))
		assert.Equal(t, status, result.Status)
		assert.NoError(t, result.Err)
	})

	t.Run("status retry with backoff", func(t *testing.T) {
		body := []byte(`{"message":"hi"}`)
		status := 429
		url := "url"
		mockClient := mocks.NewMockAPIClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(body, status, nil)

		logger := zap.NewNop().Sugar()
		timeout := 100 * time.Millisecond

		client := v2.NewClientFromRaw(logger, mockClient, timeout)
		result, retryFunc := client.MakeRetryFunc(url)

		state, err := retryFunc(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, retry.Backoff, state)

		assert.Equal(t, string(body), string(result.Body))
		assert.Equal(t, status, result.Status)
		assert.NoError(t, result.Err)
	})

	t.Run("status fail", func(t *testing.T) {
		body := []byte(`{"message":"hi"}`)
		status := 1
		url := "url"
		mockClient := mocks.NewMockAPIClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(body, status, nil)

		logger := zap.NewNop().Sugar()
		timeout := 100 * time.Millisecond

		client := v2.NewClientFromRaw(logger, mockClient, timeout)
		result, retryFunc := client.MakeRetryFunc(url)

		state, err := retryFunc(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, retry.Fail, state)

		assert.Equal(t, string(body), string(result.Body))
		assert.Equal(t, status, result.Status)
		assert.NoError(t, result.Err)
	})

	t.Run("status fall through", func(t *testing.T) {
		body := []byte(`{"message":"hi"}`)
		status := 0
		url := "url"
		mockClient := mocks.NewMockAPIClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(body, status, nil)

		logger := zap.NewNop().Sugar()
		timeout := 100 * time.Millisecond

		client := v2.NewClientFromRaw(logger, mockClient, timeout)
		result, retryFunc := client.MakeRetryFunc(url)

		state, err := retryFunc(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, retry.Fail, state)

		assert.Equal(t, string(body), string(result.Body))
		assert.Equal(t, status, result.Status)
		assert.NoError(t, result.Err)
	})
}
