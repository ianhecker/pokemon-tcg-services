package cards_test

// import (
// 	"os"
// 	"path/filepath"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"

// 	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg/v1/cards"
// 	"github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"
// )

// func readTestdata(t *testing.T, name string) []byte {
// 	path := filepath.Join("testdata", name)
// 	bytes, err := os.ReadFile(path)
// 	if err != nil {
// 		assert.Failf(t, "reading testdata: %s: %v", path, err)
// 	}
// 	return bytes
// }

// func TestResponse_UnmarshalJSON(t *testing.T) {
// 	t.Run("happy path", func(t *testing.T) {
// 		want := cards.Card{
// 			ID:          "pokemon-xy-evolutions-venusaur-ex-ultra-rare",
// 			Name:        "Venusaur EX",
// 			Number:      "1/108",
// 			Rarity:      "Ultra Rare",
// 			Set:         "XY - Evolutions",
// 			TCGPlayerID: cards.CardID("124014"),
// 			Pricing: []cards.ConditionPricing{
// 				cards.ConditionPricing{
// 					Condition:   pokemontcg.Damaged,
// 					LastUpdated: time.Unix(1756329181, 0).UTC(),
// 					Market:      1.05,
// 					MaxPrice30d: 1.05,
// 					MinPrice30d: 1.04,
// 				},
// 			},
// 		}
// 		bytes := readTestdata(t, "card.json")

// 		var response cards.PricingResponse
// 		err := response.UnmarshalJSON(bytes)
// 		require.NoError(t, err)

// 		card, err := response.GetCardIndex(0)
// 		assert.NoError(t, err)
// 		assert.Equal(t, want, card)
// 	})

// 	t.Run("bad json", func(t *testing.T) {
// 		bytes := readTestdata(t, "malformed_card.json")

// 		var response cards.PricingResponse
// 		err := response.UnmarshalJSON(bytes)
// 		require.ErrorContains(t, err, "invalid character ',' after top-level value")
// 	})

// 	t.Run("empty cards", func(t *testing.T) {
// 		bytes := readTestdata(t, "empty_cards.json")

// 		var response cards.PricingResponse
// 		err := response.UnmarshalJSON(bytes)
// 		require.ErrorContains(t, err, "zero cards in response")
// 	})
// }
