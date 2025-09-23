package cards_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg/v1/cards"
	"github.com/stretchr/testify/assert"
)

func readTestdata(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name)
	bytes, err := os.ReadFile(path)
	if err != nil {
		assert.Failf(t, "reading testdata: %s: %v", path, err)
	}
	return bytes
}

func TestCard_UnmarshalJSON(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		bytes := readTestdata(t, "card.json")

		expected := cards.Card{
			ID:          "pokemon-xy-evolutions-venusaur-ex-ultra-rare",
			TCGPlayerID: "124014",
			Name:        "Venusaur EX",
			Number:      "1/108",
			Rarity:      "Ultra Rare",
			Set:         "XY - Evolutions",
			Prices: cards.PrintingConditionPrices{
				"Holofoil": cards.ConditionPrices{
					NearMint: cards.Prices{
						Market:      3.25,
						MinPrice30d: 3.2,
						MaxPrice30d: 3.49,
						LastUpdated: time.Unix(1756380625, 0).UTC(),
					},
					LightlyPlayed: cards.Prices{
						Market:      3.01,
						MinPrice30d: 2.9,
						MaxPrice30d: 3.01,
						LastUpdated: time.Unix(1756380625, 0).UTC(),
					},
					ModeratelyPlayed: cards.Prices{
						Market:      2.31,
						MinPrice30d: 2.25,
						MaxPrice30d: 2.31,
						LastUpdated: time.Unix(1756380625, 0).UTC(),
					},
					HeavilyPlayed: cards.Prices{
						Market:      1.56,
						MinPrice30d: 1.5,
						MaxPrice30d: 1.56,
						LastUpdated: time.Unix(1756380625, 0).UTC(),
					},
					Damaged: cards.Prices{
						Market:      1.06,
						MinPrice30d: 1.04,
						MaxPrice30d: 1.06,
						LastUpdated: time.Unix(1756380625, 0).UTC(),
					},
				},
			},
		}
		var card cards.Card
		err := card.UnmarshalJSON(bytes)
		assert.NoError(t, err)
		assert.Equal(t, expected, card)
	})

	t.Run("bad json", func(t *testing.T) {
		bytes := readTestdata(t, "bad.json")

		var card cards.Card
		err := card.UnmarshalJSON(bytes)
		assert.ErrorContains(t, err, "error unmarshaling card: invalid character")
	})
}
