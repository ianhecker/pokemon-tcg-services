package justtcg_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/ianhecker/pokemon-tcg-services/internal/config"
	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg"
	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg/v1/cards"
	"github.com/ianhecker/pokemon-tcg-services/internal/retry"
	"github.com/ianhecker/pokemon-tcg-services/internal/testkit/fixtures"
	"github.com/ianhecker/pokemon-tcg-services/internal/testkit/generate"

	mocks "github.com/ianhecker/pokemon-tcg-services/internal/mocks/networkingmocks"
)

func TestClient_NewClient(t *testing.T) {
	logger := zap.NewNop().Sugar()
	token := config.MakeToken("token")
	client := justtcg.NewClient(logger, token)

	assert.Implements(t, (*justtcg.ClientInterface)(nil), client)
}

func TestClient_GetPricing(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		expected := generate.DefaultCard()
		ID := expected.TCGPlayerID

		body := fixtures.Read(t, "response.json")
		status := 200
		url := cards.GetCardByID(ID)
		mockClient := mocks.NewMockClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(body, status, nil)

		logger := zap.NewNop().Sugar()
		client := justtcg.NewClientFromRaw(logger, mockClient)

		card, err := client.GetPricing(context.Background(), ID)
		require.NoError(t, err)
		assert.Equal(t, expected, card)
	})
	t.Run("retryable error", func(t *testing.T) {
		expected := generate.DefaultCard()
		ID := expected.TCGPlayerID

		url := cards.GetCardByID(ID)
		mockClient := mocks.NewMockClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(nil, 0, errors.New("error"))

		logger := zap.NewNop().Sugar()
		client := justtcg.NewClientFromRaw(logger, mockClient)

		_, err := client.GetPricing(context.Background(), ID)
		assert.ErrorContains(t, err, "client retrying")
	})
	t.Run("decode error", func(t *testing.T) {
		body := []byte(`bad json`)
		status := 200
		url := cards.GetCardByID("")
		mockClient := mocks.NewMockClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(body, status, nil)

		logger := zap.NewNop().Sugar()
		client := justtcg.NewClientFromRaw(logger, mockClient)

		_, err := client.GetPricing(context.Background(), "")
		assert.ErrorContains(t, err, "client decoding")
	})
	t.Run("map error", func(t *testing.T) {
		expected := generate.DefaultCard()
		ID := expected.TCGPlayerID

		body := []byte(`{"data":[{"tcgplayerid":""}]}`)
		status := 200
		url := cards.GetCardByID(ID)
		mockClient := mocks.NewMockClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(body, status, nil)

		logger := zap.NewNop().Sugar()
		client := justtcg.NewClientFromRaw(logger, mockClient)

		_, err := client.GetPricing(context.Background(), ID)
		assert.ErrorContains(t, err, "client mapping")
	})
}

func TestClient_MakeRetryFunc(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		body := []byte(`{"message":"hi"}`)
		status := 200
		url := "url"
		mockClient := mocks.NewMockClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(body, status, nil)

		logger := zap.NewNop().Sugar()
		client := justtcg.NewClientFromRaw(logger, mockClient)
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
		mockClient := mocks.NewMockClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(nil, 0, expected)

		logger := zap.NewNop().Sugar()
		client := justtcg.NewClientFromRaw(logger, mockClient)
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
		mockClient := mocks.NewMockClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(body, status, nil)

		logger := zap.NewNop().Sugar()
		client := justtcg.NewClientFromRaw(logger, mockClient)
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
		mockClient := mocks.NewMockClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(body, status, nil)

		logger := zap.NewNop().Sugar()
		client := justtcg.NewClientFromRaw(logger, mockClient)
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
		mockClient := mocks.NewMockClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(body, status, nil)

		logger := zap.NewNop().Sugar()
		client := justtcg.NewClientFromRaw(logger, mockClient)
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
		mockClient := mocks.NewMockClientInterface(t)
		mockClient.
			On("Get", mock.Anything, url).
			Return(body, status, nil)

		logger := zap.NewNop().Sugar()
		client := justtcg.NewClientFromRaw(logger, mockClient)
		result, retryFunc := client.MakeRetryFunc(url)

		state, err := retryFunc(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, retry.Fail, state)

		assert.Equal(t, string(body), string(result.Body))
		assert.Equal(t, status, result.Status)
		assert.NoError(t, result.Err)
	})
}
