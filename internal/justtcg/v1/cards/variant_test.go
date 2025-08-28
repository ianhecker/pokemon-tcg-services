package cards_test

import (
	"testing"
	"time"

	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg/v1/cards"
	"github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVariants_ToMap(t *testing.T) {
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
		expected := map[pokemontcg.Condition]cards.Prices{
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
		m, err := variants.ToMap()
		require.NoError(t, err)

		price, ok := m[pokemontcg.NearMint]
		require.True(t, ok)
		assert.Equal(t, expected[pokemontcg.NearMint], price)

		price, ok = m[pokemontcg.LightlyPlayed]
		require.True(t, ok)
		assert.Equal(t, expected[pokemontcg.LightlyPlayed], price)

		price, ok = m[pokemontcg.ModeratelyPlayed]
		require.True(t, ok)
		assert.Equal(t, expected[pokemontcg.ModeratelyPlayed], price)

		price, ok = m[pokemontcg.HeavilyPlayed]
		require.True(t, ok)
		assert.Equal(t, expected[pokemontcg.HeavilyPlayed], price)

		price, ok = m[pokemontcg.Damaged]
		require.True(t, ok)
		assert.Equal(t, expected[pokemontcg.Damaged], price)
	})

	t.Run("duplicate condition", func(t *testing.T) {
		variants := cards.Variants{
			cards.Variant{
				Condition:   pokemontcg.NearMint,
				Price:       1.0,
				MinPrice30d: 1.1,
				MaxPrice30d: 1.2,
				LastUpdated: 1756380001,
			},
			cards.Variant{
				Condition:   pokemontcg.NearMint,
				Price:       2.0,
				MinPrice30d: 2.1,
				MaxPrice30d: 2.2,
				LastUpdated: 1756380002,
			},
		}

		_, err := variants.ToMap()
		assert.ErrorContains(t, err, "duplicated condition in prices: 'Near Mint'")
	})
}
