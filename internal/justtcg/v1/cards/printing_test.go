package cards_test

import (
	"testing"
	"time"

	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg/v1/cards"
	"github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintingConditionPrices_MakePrintingConditionPricesFromVariants(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		variants := []cards.Variant{
			cards.Variant{
				Printing:    "Holofoil",
				Condition:   pokemontcg.NearMint,
				Price:       1.0,
				MinPrice30d: 1.1,
				MaxPrice30d: 1.2,
				LastUpdated: 1756380001,
			},
			cards.Variant{
				Printing:    "Holofoil",
				Condition:   pokemontcg.LightlyPlayed,
				Price:       2.0,
				MinPrice30d: 2.1,
				MaxPrice30d: 2.2,
				LastUpdated: 1756380002,
			},
			cards.Variant{
				Printing:    "Holofoil",
				Condition:   pokemontcg.ModeratelyPlayed,
				Price:       3.0,
				MinPrice30d: 3.1,
				MaxPrice30d: 3.2,
				LastUpdated: 1756380003,
			},
			cards.Variant{
				Printing:    "Unlimited Holofoil",
				Condition:   pokemontcg.HeavilyPlayed,
				Price:       4.1,
				MinPrice30d: 4.2,
				MaxPrice30d: 4.3,
				LastUpdated: 1756380004,
			},
			cards.Variant{
				Printing:    "Unlimited Holofoil",
				Condition:   pokemontcg.Damaged,
				Price:       5.1,
				MinPrice30d: 5.2,
				MaxPrice30d: 5.3,
				LastUpdated: 1756380005,
			},
		}
		expected := map[cards.Printing]cards.ConditionPrices{
			"Holofoil": cards.ConditionPrices{
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
			},
			"Unlimited Holofoil": cards.ConditionPrices{
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
			},
		}
		m := cards.MakePrintingConditionPricesFromVariants(variants)

		prices, ok := m["Holofoil"]
		require.True(t, ok)
		assert.Equal(t, expected["Holofoil"], prices)

		prices, ok = m["Holofoil"]
		require.True(t, ok)
		assert.Equal(t, expected["Holofoil"], prices)

		prices, ok = m["Holofoil"]
		require.True(t, ok)
		assert.Equal(t, expected["Holofoil"], prices)

		prices, ok = m["Unlimited Holofoil"]
		require.True(t, ok)
		assert.Equal(t, expected["Unlimited Holofoil"], prices)

		prices, ok = m["Unlimited Holofoil"]
		require.True(t, ok)
		assert.Equal(t, expected["Unlimited Holofoil"], prices)
	})
}
