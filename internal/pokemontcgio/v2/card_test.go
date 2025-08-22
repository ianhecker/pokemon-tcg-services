package v2_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	v2 "github.com/ianhecker/pokemon-tcg-services/internal/pokemontcgio/v2"
)

func TestCardByID(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		card := "my-card-123"
		api := v2.CardByID(card)
		assert.Equal(t, api, fmt.Sprintf("api.pokemontcg.io/v2/cards/%s", card))
	})
}
