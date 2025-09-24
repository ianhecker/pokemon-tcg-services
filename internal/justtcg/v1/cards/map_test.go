package cards_test

import (
	"os"
	"path/filepath"
	"testing"

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

	t.Run("empty data", func(t *testing.T) {
		responseDTO := readTestData(t, "/validate/empty/data.json")

		err := cards.Validate(responseDTO)
		assert.ErrorContains(t, err, "data field is empty")
	})

	t.Run("empty TCGPlayer ID", func(t *testing.T) {
		responseDTO := readTestData(t, "/validate/empty/tcgplayerid.json")

		err := cards.Validate(responseDTO)
		assert.ErrorContains(t, err, "card: TCGPlayerId is empty")
	})

	t.Run("empty variants", func(t *testing.T) {
		responseDTO := readTestData(t, "/validate/empty/variants.json")

		err := cards.Validate(responseDTO)
		assert.ErrorContains(t, err, "card: variants are empty")
	})

	t.Run("empty printing", func(t *testing.T) {
		responseDTO := readTestData(t, "/validate/empty/printing.json")

		err := cards.Validate(responseDTO)
		assert.ErrorContains(t, err, "card: variants: printing is empty")
	})

	t.Run("empty condition", func(t *testing.T) {
		responseDTO := readTestData(t, "/validate/empty/condition.json")

		err := cards.Validate(responseDTO)
		assert.ErrorContains(t, err, "card: variants: condition is empty")
	})
}
