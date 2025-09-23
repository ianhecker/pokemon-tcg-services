package cards_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg/v1/cards"
	"github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"
)

func TestConditionPrices_MakeConditionPricesFromMap(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		m := map[pokemontcg.Condition]cards.Prices{
			pokemontcg.NearMint: cards.Prices{
				Market:      1.0,
				MinPrice30d: 1.1,
				MaxPrice30d: 1.2,
				LastUpdated: time.Unix(1756380001, 0).UTC(),
			},
			pokemontcg.LightlyPlayed: cards.Prices{
				Market:      2.0,
				MinPrice30d: 2.1,
				MaxPrice30d: 2.2,
				LastUpdated: time.Unix(1756380002, 0).UTC(),
			},
			pokemontcg.ModeratelyPlayed: cards.Prices{
				Market:      3.0,
				MinPrice30d: 3.1,
				MaxPrice30d: 3.2,
				LastUpdated: time.Unix(1756380003, 0).UTC(),
			},
			pokemontcg.HeavilyPlayed: cards.Prices{
				Market:      4.1,
				MinPrice30d: 4.2,
				MaxPrice30d: 4.3,
				LastUpdated: time.Unix(1756380004, 0).UTC(),
			},
			pokemontcg.Damaged: cards.Prices{
				Market:      5.1,
				MinPrice30d: 5.2,
				MaxPrice30d: 5.3,
				LastUpdated: time.Unix(1756380005, 0).UTC(),
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
		conditionPrices := cards.MakeConditionPricesFromMap(m)

		assert.Equal(t, expected.NearMint, conditionPrices.NearMint)
		assert.Equal(t, expected.LightlyPlayed, conditionPrices.LightlyPlayed)
		assert.Equal(t, expected.ModeratelyPlayed, conditionPrices.ModeratelyPlayed)
		assert.Equal(t, expected.HeavilyPlayed, conditionPrices.HeavilyPlayed)
		assert.Equal(t, expected.Damaged, conditionPrices.Damaged)
	})
}
