package proxy_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/ianhecker/pokemon-tcg-services/internal/mocks/proxymocks"
	"github.com/ianhecker/pokemon-tcg-services/internal/networking/proxy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func URL(t *testing.T, s string) *url.URL {
	url, err := url.Parse(s)
	require.NoError(t, err)
	return url
}

type errorBody struct{}

func (e *errorBody) Read(p []byte) (int, error) {
	return 0, errors.New("error")
}

func (e *errorBody) Close() error {
	return nil
}

func TestProxy_Do(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		url := URL(t, "www.ianhecker.com")
		req := &http.Request{URL: url}

		resp := &http.Response{
			StatusCode: 200,
			Header:     make(http.Header),
		}
		data := `{"message":"hello world"}`
		resp.Body = io.NopCloser(bytes.NewReader([]byte(data)))
		resp.ContentLength = int64(len(data))

		client := proxymocks.NewMockClientInterface(t)
		client.
			On("Do", req).
			Return(resp, nil)

		request := proxy.Request{Req: req}

		proxy := proxy.NewProxyFromRaw(client)
		response := proxy.Do(request)

		assert.NoError(t, response.Err)
		assert.Equal(t, 200, response.Status)
		assert.Equal(t, data, string(response.Body))
		client.AssertExpectations(t)
	})
	t.Run("nil response", func(t *testing.T) {
		url := URL(t, "www.ianhecker.com")
		req := &http.Request{URL: url}

		client := proxymocks.NewMockClientInterface(t)
		client.
			On("Do", req).
			Return(nil, nil)

		request := proxy.Request{Req: req}

		proxy := proxy.NewProxyFromRaw(client)
		response := proxy.Do(request)

		assert.ErrorContains(t, response.Err, "nil response")
		client.AssertExpectations(t)
	})
	t.Run("response err", func(t *testing.T) {
		url := URL(t, "www.ianhecker.com")
		req := &http.Request{URL: url}

		resp := &http.Response{
			StatusCode: 404,
		}

		client := proxymocks.NewMockClientInterface(t)
		client.
			On("Do", req).
			Return(resp, errors.New("error"))

		request := proxy.Request{Req: req}

		proxy := proxy.NewProxyFromRaw(client)
		response := proxy.Do(request)

		assert.ErrorContains(t, response.Err, "proxy: do: error")
		assert.Equal(t, 404, response.Status)
		client.AssertExpectations(t)
	})
	t.Run("read body error", func(t *testing.T) {
		url := URL(t, "www.ianhecker.com")
		req := &http.Request{URL: url}

		resp := &http.Response{
			StatusCode: 200,
			Body:       &errorBody{},
		}

		client := proxymocks.NewMockClientInterface(t)
		client.
			On("Do", req).
			Return(resp, nil)

		request := proxy.Request{Req: req}

		proxy := proxy.NewProxyFromRaw(client)
		response := proxy.Do(request)

		assert.ErrorContains(t, response.Err, "error")
		assert.Equal(t, 200, response.Status)
		client.AssertExpectations(t)
	})
}
