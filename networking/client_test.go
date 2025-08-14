package networking_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ianhecker/pokemon-tcg-services/internal/mocks"
	"github.com/ianhecker/pokemon-tcg-services/networking"
)

func TestClient_NewRequest(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		method := http.MethodGet
		url := "https://google.com"

		client := networking.NewClientFromRaw(httpClientMock)
		request, err := client.NewRequest(method, url)

		assert.NoError(t, err)
		assert.Equal(t, method, request.Method())

		gotURL, err := request.URL()
		assert.NoError(t, err)
		assert.Equal(t, url, gotURL)

		httpClientMock.AssertExpectations(t)
	})

	t.Run("bad method", func(t *testing.T) {
		httpClientMock := mocks.NewMockHTTPClient(t)

		client := networking.NewClientFromRaw(httpClientMock)
		_, err := client.NewRequest("not a method", "")

		assert.ErrorContains(t, err, "invalid method \"not a method\"")

		httpClientMock.AssertExpectations(t)
	})
}
