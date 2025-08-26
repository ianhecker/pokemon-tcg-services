package pokemonpricetracker_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/ianhecker/pokemon-tcg-services/internal/mocks"
	"github.com/ianhecker/pokemon-tcg-services/internal/pokemonpricetracker"
	"github.com/ianhecker/pokemon-tcg-services/internal/retry"
)

func TestClient_MakeRetryFunc(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		body := []byte(`{"message":"hi"}`)
		status := 200
		url := "url"
		mockClient := mocks.NewMockHttpClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(body, status, nil)

		logger := zap.NewNop().Sugar()
		client := pokemonpricetracker.NewClientFromRaw(logger, mockClient)
		result, retryFunc := client.MakeRetryFunc(url)

		state, err := retryFunc(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, retry.Complete, state)

		assert.Equal(t, string(body), string(result.Body))
		assert.Equal(t, status, result.Status)
		assert.NoError(t, result.Err)
	})

	t.Run("http client error", func(t *testing.T) {
		url := "url"
		expected := errors.New("error")
		mockClient := mocks.NewMockHttpClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(nil, 0, expected)

		logger := zap.NewNop().Sugar()
		client := pokemonpricetracker.NewClientFromRaw(logger, mockClient)
		result, retryFunc := client.MakeRetryFunc(url)

		state, err := retryFunc(context.Background())
		assert.ErrorIs(t, err, expected)
		assert.Equal(t, retry.Fail, state)
		assert.ErrorIs(t, result.Err, expected)
	})

	t.Run("status retry", func(t *testing.T) {
		body := []byte(`{"message":"hi"}`)
		status := 504
		url := "url"
		mockClient := mocks.NewMockHttpClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(body, status, nil)

		logger := zap.NewNop().Sugar()
		client := pokemonpricetracker.NewClientFromRaw(logger, mockClient)
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
		mockClient := mocks.NewMockHttpClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(body, status, nil)

		logger := zap.NewNop().Sugar()
		client := pokemonpricetracker.NewClientFromRaw(logger, mockClient)
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
		mockClient := mocks.NewMockHttpClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(body, status, nil)

		logger := zap.NewNop().Sugar()
		client := pokemonpricetracker.NewClientFromRaw(logger, mockClient)
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
		mockClient := mocks.NewMockHttpClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(body, status, nil)

		logger := zap.NewNop().Sugar()
		client := pokemonpricetracker.NewClientFromRaw(logger, mockClient)
		result, retryFunc := client.MakeRetryFunc(url)

		state, err := retryFunc(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, retry.Fail, state)

		assert.Equal(t, string(body), string(result.Body))
		assert.Equal(t, status, result.Status)
		assert.NoError(t, result.Err)
	})
}
