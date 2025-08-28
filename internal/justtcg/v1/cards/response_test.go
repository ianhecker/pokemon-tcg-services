package cards_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg/v1/cards"
)

func TestResponse_UnmarshalJSON(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		expected := cards.Card{
			ID:          "pokemon-xy-evolutions-venusaur-ex-ultra-rare",
			TCGPlayerID: "124014",
			Name:        "Venusaur EX",
			Number:      "1/108",
			Rarity:      "Ultra Rare",
			Set:         "XY - Evolutions",
			Prices: cards.ConditionPrices{
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
		}
		bytes := readTestdata(t, "response.json")

		var response cards.Response
		err := response.UnmarshalJSON(bytes)
		require.NoError(t, err)

		card, err := response.GetCardIndex(0)
		assert.NoError(t, err)
		assert.Equal(t, expected, card)
	})

	t.Run("bad json", func(t *testing.T) {
		bytes := readTestdata(t, "bad.json")

		var response cards.Response
		err := response.UnmarshalJSON(bytes)
		require.ErrorContains(t, err, "error unmarshaling response")
	})

	t.Run("empty cards", func(t *testing.T) {
		bytes := readTestdata(t, "empty_cards.json")

		var response cards.Response
		err := response.UnmarshalJSON(bytes)
		require.ErrorContains(t, err, "zero cards in response")
	})
}
