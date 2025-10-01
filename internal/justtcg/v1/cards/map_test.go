package cards_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg/v1/cards"
	"github.com/ianhecker/pokemon-tcg-services/internal/testkit"
	"github.com/ianhecker/pokemon-tcg-services/internal/testkit/generate"
)

func TestValidate(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		responseDTO := generate.DefaultResponseDTO()

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
		responseDTO := generate.DefaultResponseDTO()
		expected := generate.DefaultCard()

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

				bytes := testkit.ReadTestdata(t, test.filepath)

				responseDTO, err := cards.Decode(bytes)
				require.NoError(t, err)

				_, err = cards.Map(responseDTO)
				assert.ErrorContains(t, err, test.err)
			})
		}
	})
}
