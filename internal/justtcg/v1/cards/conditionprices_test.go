package cards_test

import (
	"testing"
	"time"

	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg/v1/cards"
	"github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConditionPrices_MakeConditionFromVariants(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		variants := cards.Variants{
			cards.Variant{
				Condition:   pokemontcg.NearMint,
				Price:       1.0,
				MinPrice30d: 1.1,
				MaxPrice30d: 1.2,
				LastUpdated: 1756380001,
			},
			cards.Variant{
				Condition:   pokemontcg.LightlyPlayed,
				Price:       2.0,
				MinPrice30d: 2.1,
				MaxPrice30d: 2.2,
				LastUpdated: 1756380002,
			},
			cards.Variant{
				Condition:   pokemontcg.ModeratelyPlayed,
				Price:       3.0,
				MinPrice30d: 3.1,
				MaxPrice30d: 3.2,
				LastUpdated: 1756380003,
			},
			cards.Variant{
				Condition:   pokemontcg.HeavilyPlayed,
				Price:       4.1,
				MinPrice30d: 4.2,
				MaxPrice30d: 4.3,
				LastUpdated: 1756380004,
			},
			cards.Variant{
				Condition:   pokemontcg.Damaged,
				Price:       5.1,
				MinPrice30d: 5.2,
				MaxPrice30d: 5.3,
				LastUpdated: 1756380005,
			},
		}
		expected := cards.ConditionPrices{
			NearMint: cards.Prices{
				Market:      1.0,
				MinPrice30d: 1.1,
				MaxPrice30d: 1.2,
				LastUpdated: time.Unix(1756380001, 0).UTC(),
			},
			LightlyPlayed: cards.Prices{
				Market:      2.0,
				MinPrice30d: 2.1,
				MaxPrice30d: 2.2,
				LastUpdated: time.Unix(1756380002, 0).UTC(),
			},
			ModeratelyPlayed: cards.Prices{
				Market:      3.0,
				MinPrice30d: 3.1,
				MaxPrice30d: 3.2,
				LastUpdated: time.Unix(1756380003, 0).UTC(),
			},
			HeavilyPlayed: cards.Prices{
				Market:      4.1,
				MinPrice30d: 4.2,
				MaxPrice30d: 4.3,
				LastUpdated: time.Unix(1756380004, 0).UTC(),
			},
			Damaged: cards.Prices{
				Market:      5.1,
				MinPrice30d: 5.2,
				MaxPrice30d: 5.3,
				LastUpdated: time.Unix(1756380005, 0).UTC(),
			},
		}
		m, err := variants.ToMap()
		require.NoError(t, err)

		conditionPrices, err := cards.MakeConditionPricesFromVariants(m)
		require.NoError(t, err)

		assert.Equal(t, expected.NearMint, conditionPrices.NearMint)
		assert.Equal(t, expected.LightlyPlayed, conditionPrices.LightlyPlayed)
		assert.Equal(t, expected.ModeratelyPlayed, conditionPrices.ModeratelyPlayed)
		assert.Equal(t, expected.HeavilyPlayed, conditionPrices.HeavilyPlayed)
		assert.Equal(t, expected.Damaged, conditionPrices.Damaged)
	})

	t.Run("at least price set", func(t *testing.T) {
		variants := cards.Variants{
			cards.Variant{
				Condition:   pokemontcg.NearMint,
				Price:       1.0,
				MinPrice30d: 1.1,
				MaxPrice30d: 1.2,
				LastUpdated: 1756380001,
			},
		}
		expected := cards.ConditionPrices{
			NearMint: cards.Prices{
				Market:      1.0,
				MinPrice30d: 1.1,
				MaxPrice30d: 1.2,
				LastUpdated: time.Unix(1756380001, 0).UTC(),
			},
		}
		m, err := variants.ToMap()
		require.NoError(t, err)

		conditionPrices, err := cards.MakeConditionPricesFromVariants(m)
		require.NoError(t, err)

		assert.Equal(t, expected.NearMint, conditionPrices.NearMint)
		assert.Equal(t, cards.Prices{}, conditionPrices.LightlyPlayed)
		assert.Equal(t, cards.Prices{}, conditionPrices.ModeratelyPlayed)
		assert.Equal(t, cards.Prices{}, conditionPrices.HeavilyPlayed)
		assert.Equal(t, cards.Prices{}, conditionPrices.Damaged)
	})

	t.Run("no prices given", func(t *testing.T) {
		variants := cards.Variants{}
		m, err := variants.ToMap()
		require.NoError(t, err)

		_, err = cards.MakeConditionPricesFromVariants(m)
		require.ErrorContains(t, err, "given zero prices")
	})
}
