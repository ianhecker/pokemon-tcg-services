package client_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ianhecker/pokemon-tcg-services/client"
)

func mockResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		Status:     http.StatusText(statusCode),
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}
}

func TestResponse_ReadBody(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		expected := `{"message":"hello"}`
		resp := mockResponse(200, expected)
		response := client.NewResponse(resp)

		body, err := response.ReadBody()
		assert.NoError(t, err)
		assert.Equal(t, expected, string(body))
	})

	t.Run("nil request", func(t *testing.T) {
		response := client.NewResponse(nil)

		_, err := response.ReadBody()
		assert.ErrorContains(t, err, "response is nil")
	})
}
