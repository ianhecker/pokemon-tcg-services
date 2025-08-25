package pokemontcg_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func readTestdata(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name)
	bytes, err := os.ReadFile(path)
	if err != nil {
		assert.Failf(t, "reading testdata: %s: %v", path, err)
	}
	return bytes
}

func TestPricingResponse_Unmarshal(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		want := pokemontcg.Card{
			ID:     "xy1-1",
			Name:   "Venusaur-EX",
			Number: 1,
			Set: pokemontcg.Set{
				ID:          "xy1",
				Name:        "XY",
				Series:      "XY",
				Total:       146,
				ReleaseDate: "2014/02/05",
			},
			LastUpdated: "2025-08-25T17:16:54.954Z",
		}
		bytes := readTestdata(t, "card.json")

		var response pokemontcg.PricingResponse
		err := response.UnmarshalJSON(bytes)
		require.NoError(t, err)

		card, err := response.GetCardIndex(0)
		assert.NoError(t, err)
		assert.Equal(t, want, card)
	})

	t.Run("bad json", func(t *testing.T) {
		bytes := readTestdata(t, "malformed_card.json")

		var response pokemontcg.PricingResponse
		err := response.UnmarshalJSON(bytes)
		require.ErrorContains(t, err, "invalid character ',' after top-level value")
	})

	t.Run("empty cards", func(t *testing.T) {
		bytes := readTestdata(t, "empty_cards.json")

		var response pokemontcg.PricingResponse
		err := response.UnmarshalJSON(bytes)
		require.ErrorContains(t, err, "zero cards in response")
	})
}
