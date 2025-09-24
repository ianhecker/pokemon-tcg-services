package cards_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg/v1/cards"
)

func readTestData(t *testing.T, filename string) cards.ResponseDTO {
	path := filepath.Join("testdata", filename)
	bytes, err := os.ReadFile(path)
	require.NoError(t, err)

	responseDTO, err := cards.Decode(bytes, cards.UseNumber())
	require.NoError(t, err)
	return responseDTO
}

func TestValidate(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		responseDTO := readTestData(t, "response.json")

		err := cards.Validate(responseDTO)
		assert.NoError(t, err)
	})
	t.Run("errors", func(t *testing.T) {
		tests := []struct {
			name     string
			response cards.ResponseDTO
			err      string
		}{
			{"empty data", cards.ResponseDTO{Data: nil}, "data field is empty"},
			{"empty TCGPlayerID", cards.ResponseDTO{Data: []cards.CardDTO{{TCGPlayerID: ""}}}, "TCGPlayerId is empty"},
			{"empty variants", cards.ResponseDTO{Data: []cards.CardDTO{{TCGPlayerID: "106999", Variants: nil}}}, "variants are empty"},
			{"empty printing", cards.ResponseDTO{Data: []cards.CardDTO{{TCGPlayerID: "106999", Variants: []cards.VariantDTO{{Printing: ""}}}}}, "printing is empty"},
			{"empty condition", cards.ResponseDTO{Data: []cards.CardDTO{{TCGPlayerID: "106999", Variants: []cards.VariantDTO{{Printing: "1st Edition Holofoil", Condition: ""}}}}}, "condition is empty"},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {

				err := cards.Validate(test.response)
				assert.ErrorContains(t, err, test.err)
			})
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		responseDTO := readTestData(t, "response.json")

		got, err := cards.Map(responseDTO)
		require.NoError(t, err)

		equal := cmp.Equal(expected, got)
		assert.True(t, equal, cmp.Diff(expected, got))
	})

	t.Run("errors", func(t *testing.T) {
		tests := []struct {
			name     string
			filepath string
			err      string
		}{
			{"validation error", "bad/data.json", "error validating response"},
			{"TCGPlayerID error", "bad/tcgplayerid.json", "error parsing TCGPlayerID"},
			{"variant error", "bad/variants.json", "error mapping variants"},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {

				responseDTO := readTestData(t, test.filepath)
				_, err := cards.Map(responseDTO)
				assert.ErrorContains(t, err, test.err)
			})
		}
	})
}

func fptr(f float64) *cards.Price {
	price := cards.Price(f)
	return &price
}

func tptr(t time.Time) *time.Time {
	return &t
}

var expected = cards.Card{
	ID:          "pokemon-base-set-shadowless-charizard-holo-rare",
	TCGPlayerID: "106999",
	Name:        "Charizard",
	Number:      "004/102",
	Rarity:      "Holo Rare",
	Set:         "Base Set (Shadowless)",
	Prices: cards.Printings{
		"1st Edition Holofoil": cards.Conditions{
			NearMint: &cards.Prices{
				Market:      fptr(503.49),
				MinPrice30d: fptr(500.99),
				MaxPrice30d: fptr(503.49),
				LastUpdated: tptr(time.Unix(1758706033, 0).UTC()),
			},
			ModeratelyPlayed: &cards.Prices{
				Market:      fptr(5495),
				MinPrice30d: fptr(5495),
				MaxPrice30d: fptr(5495),
				LastUpdated: tptr(time.Unix(1754787579, 0).UTC()),
			},
			HeavilyPlayed: &cards.Prices{
				Market:      cards.MakePrice(225),
				LastUpdated: tptr(time.Unix(1743656485, 0).UTC()),
			},
		},
		"Unlimited Holofoil": cards.Conditions{
			NearMint: &cards.Prices{
				Market:      fptr(1923.19),
				MinPrice30d: fptr(1700),
				MaxPrice30d: fptr(1923.19),
				LastUpdated: tptr(time.Unix(1758706033, 0).UTC()),
			},
			LightlyPlayed: &cards.Prices{
				Market:      fptr(1299),
				MinPrice30d: fptr(1299),
				MaxPrice30d: fptr(1299),
				LastUpdated: tptr(time.Unix(1758706033, 0).UTC()),
			},
			ModeratelyPlayed: &cards.Prices{
				Market:      fptr(764.78),
				MinPrice30d: fptr(761.02),
				MaxPrice30d: fptr(813.73),
				LastUpdated: tptr(time.Unix(1758706033, 0).UTC()),
			},
			HeavilyPlayed: &cards.Prices{
				Market:      fptr(645.05),
				MinPrice30d: fptr(642.95),
				MaxPrice30d: fptr(649),
				LastUpdated: tptr(time.Unix(1758706033, 0).UTC()),
			},
			Damaged: &cards.Prices{
				Market:      fptr(524.96),
				MinPrice30d: fptr(523.77),
				MaxPrice30d: fptr(533.29),
				LastUpdated: tptr(time.Unix(1758706033, 0).UTC()),
			},
		},
	},
}
