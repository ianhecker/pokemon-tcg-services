package cards_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg/v1/cards"
)

func TestGetCardByID(t *testing.T) {
	want := "https://api.justtcg.com/v1/cards?tcgplayerId=1234"
	got := cards.GetCardByID("1234")
	assert.Equal(t, want, got)
}
