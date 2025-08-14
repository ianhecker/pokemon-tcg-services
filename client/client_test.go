package client_test

import (
	"net/http"
	"testing"

	"github.com/ianhecker/pokemon-tcg-services/client"
	"github.com/ianhecker/pokemon-tcg-services/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_NewRequest(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		method := http.MethodGet
		url := "https://google.com"

		httpClientMock := mocks.NewMockHTTPClient(t)

		client := client.NewClientFromRaw(httpClientMock)
		request, err := client.NewRequest(method, url)

		assert.NoError(t, err)
		assert.Equal(t, method, request.Method())
		require.NotNil(t, request.URL())
		assert.Equal(t, url, request.URL().String())

		httpClientMock.AssertExpectations(t)
	})
}
