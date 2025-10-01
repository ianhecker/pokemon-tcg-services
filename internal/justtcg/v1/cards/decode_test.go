package cards_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg/v1/cards"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func readTestdata(t *testing.T, filename string) []byte {
	path := filepath.Join("testdata", filename)
	bytes, err := os.ReadFile(path)
	require.NoError(t, err)

	return bytes
}

func TestDecode(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		bytes := readTestdata(t, "response.json")
		got, err := cards.Decode(bytes)

		require.NoError(t, err)
		assert.Equal(t, expected.Response, got)
	})
	t.Run("happy path with options", func(t *testing.T) {
		called := false
		option := cards.DecodeOption(func(*json.Decoder) {
			called = true
		})
		bytes := readTestdata(t, "response.json")
		got, err := cards.Decode(bytes, option)

		require.NoError(t, err)
		assert.True(t, called)
		assert.Equal(t, expected.Response, got)
	})
	t.Run("error decoding", func(t *testing.T) {
		bytes := []byte(`bad json`)
		_, err := cards.Decode(bytes)

		assert.ErrorContains(t, err, "decode: responseDTO error")
	})
	t.Run("error extra data", func(t *testing.T) {
		bytes := []byte(`{},extra data`)
		_, err := cards.Decode(bytes)

		assert.ErrorContains(t, err, "decode: extra data after JSON")
	})
}
