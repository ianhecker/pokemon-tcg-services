package networking_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/ianhecker/pokemon-tcg-services/internal/config"
	"github.com/ianhecker/pokemon-tcg-services/internal/mocks/proxymocks"
	"github.com/ianhecker/pokemon-tcg-services/internal/networking"
	"github.com/ianhecker/pokemon-tcg-services/internal/networking/proxy"
	"github.com/ianhecker/pokemon-tcg-services/internal/testkit"
)

func TestClient_NewClient(t *testing.T) {
	logger := zap.NewNop().Sugar()
	token := config.MakeToken("token")
	client := networking.NewClient(logger, token)

	assert.Implements(t, (*networking.ClientInterface)(nil), client)
}

func TestClient_Get(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx := context.Background()
		method := proxy.GET
		url := testkit.NewURL(t, "www.google.com")

		request := proxy.Request{
			Req: &http.Request{URL: url},
			Err: nil,
		}
		token := config.MakeToken("token")
		response := proxy.Response{
			Body:   []byte(`{"message":"hello world"}`),
			Status: 200,
			Err:    nil,
		}

		proxy := proxymocks.NewMockProxyInterface(t)
		proxy.
			On("NewRequest", ctx, method, url.String()).
			Return(request, nil).
			On("SetAuthorization", request, "X-API-Key", token.Reveal()).
			Return().
			On("Do", request).
			Return(response)

		nop := zap.NewNop().Sugar()
		logger := networking.MakeLogger(nop)

		client := networking.NewClientFromRaw(logger, proxy, token)
		bytes, status, err := client.Get(ctx, url.String())

		assert.NoError(t, err)
		assert.Equal(t, bytes, response.Body)
		assert.Equal(t, status, response.Status)

		proxy.AssertExpectations(t)
	})
	t.Run("new request error", func(t *testing.T) {
		ctx := context.Background()
		method := proxy.GET
		url := testkit.NewURL(t, "www.google.com")

		request := proxy.Request{
			Req: &http.Request{URL: url},
			Err: errors.New("error"),
		}
		token := config.MakeToken("token")

		proxy := proxymocks.NewMockProxyInterface(t)
		proxy.
			On("NewRequest", ctx, method, url.String()).
			Return(request, nil)

		nop := zap.NewNop().Sugar()
		logger := networking.MakeLogger(nop)

		client := networking.NewClientFromRaw(logger, proxy, token)
		_, status, err := client.Get(ctx, url.String())

		assert.ErrorContains(t, err, "client: error")
		assert.Equal(t, status, 0)

		proxy.AssertExpectations(t)
	})
	t.Run("context canceled", func(t *testing.T) {
		ctx := context.Background()
		method := proxy.GET
		url := testkit.NewURL(t, "www.google.com")

		request := proxy.Request{
			Req: &http.Request{URL: url},
			Err: nil,
		}
		token := config.MakeToken("token")
		response := proxy.Response{
			Body:   nil,
			Status: 200,
			Err:    context.Canceled,
		}

		proxy := proxymocks.NewMockProxyInterface(t)
		proxy.
			On("NewRequest", ctx, method, url.String()).
			Return(request, nil).
			On("SetAuthorization", request, "X-API-Key", token.Reveal()).
			Return().
			On("Do", request).
			Return(response)

		nop := zap.NewNop().Sugar()
		logger := networking.MakeLogger(nop)

		client := networking.NewClientFromRaw(logger, proxy, token)
		_, status, err := client.Get(ctx, url.String())

		assert.ErrorContains(t, err, "client: context canceled")
		assert.Equal(t, status, response.Status)

		proxy.AssertExpectations(t)
	})
	t.Run("context deadline exceeded", func(t *testing.T) {
		ctx := context.Background()
		method := proxy.GET
		url := testkit.NewURL(t, "www.google.com")

		request := proxy.Request{
			Req: &http.Request{URL: url},
			Err: nil,
		}
		token := config.MakeToken("token")
		response := proxy.Response{
			Body:   nil,
			Status: 200,
			Err:    context.DeadlineExceeded,
		}

		proxy := proxymocks.NewMockProxyInterface(t)
		proxy.
			On("NewRequest", ctx, method, url.String()).
			Return(request, nil).
			On("SetAuthorization", request, "X-API-Key", token.Reveal()).
			Return().
			On("Do", request).
			Return(response)

		nop := zap.NewNop().Sugar()
		logger := networking.MakeLogger(nop)

		client := networking.NewClientFromRaw(logger, proxy, token)
		_, status, err := client.Get(ctx, url.String())

		assert.ErrorContains(t, err, "client: context deadline exceeded")
		assert.Equal(t, status, response.Status)

		proxy.AssertExpectations(t)
	})
	t.Run("other response error", func(t *testing.T) {
		ctx := context.Background()
		method := proxy.GET
		url := testkit.NewURL(t, "www.google.com")

		request := proxy.Request{
			Req: &http.Request{URL: url},
			Err: nil,
		}
		token := config.MakeToken("token")
		response := proxy.Response{
			Body:   []byte(`{"message":"hello world"}`),
			Status: 200,
			Err:    errors.New("error"),
		}

		proxy := proxymocks.NewMockProxyInterface(t)
		proxy.
			On("NewRequest", ctx, method, url.String()).
			Return(request, nil).
			On("SetAuthorization", request, "X-API-Key", token.Reveal()).
			Return().
			On("Do", request).
			Return(response)

		nop := zap.NewNop().Sugar()
		logger := networking.MakeLogger(nop)

		client := networking.NewClientFromRaw(logger, proxy, token)
		_, status, err := client.Get(ctx, url.String())

		assert.ErrorContains(t, err, "client: error")
		assert.Equal(t, status, response.Status)

		proxy.AssertExpectations(t)
	})
	t.Run("status not OK", func(t *testing.T) {
		ctx := context.Background()
		method := proxy.GET
		url := testkit.NewURL(t, "www.google.com")

		request := proxy.Request{
			Req: &http.Request{URL: url},
			Err: nil,
		}
		token := config.MakeToken("token")
		response := proxy.Response{
			Body:   []byte(`{"err":"error"}`),
			Status: 1234,
			Err:    nil,
		}

		proxy := proxymocks.NewMockProxyInterface(t)
		proxy.
			On("NewRequest", ctx, method, url.String()).
			Return(request, nil).
			On("SetAuthorization", request, "X-API-Key", token.Reveal()).
			Return().
			On("Do", request).
			Return(response)

		nop := zap.NewNop().Sugar()
		logger := networking.MakeLogger(nop)

		client := networking.NewClientFromRaw(logger, proxy, token)
		body, status, err := client.Get(ctx, url.String())

		assert.ErrorContains(t, err, "client: unexpected status code: 1234")
		assert.Equal(t, status, response.Status)
		assert.Equal(t, body, response.Body)

		proxy.AssertExpectations(t)
	})
}
