package pokemontcg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"
)

func TestCard_MakeCard(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		card, err := pokemontcg.MakeCard("base1-1")
		assert.NoError(t, err)
		assert.Equal(t, "base1-1", card.ID)
	})

	t.Run("bad ID", func(t *testing.T) {
		_, err := pokemontcg.MakeCard("bad ID")
		assert.ErrorContains(t, err, "invalid card ID: bad ID")
	})
}

func TestCardByID_SanitizeCardID(t *testing.T) {
	var tests = []struct {
		name      string
		ID        string
		sanitized string
		valid     bool
	}{
		{"happy path", "xy1-1", "xy1-1", true},
		{"happy path trim space", " xy1-1 ", "xy1-1", true},
		{"empty", "", "", false},
		{"bad ID", "bad ID", "bad ID", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sanitized, err := pokemontcg.SanitizeID(test.ID)

			if test.valid {
				assert.Nil(t, err)
				assert.Equal(t, test.sanitized, sanitized)
			} else {
				assert.ErrorContains(t, err, "invalid card ID: "+test.ID)
			}
		})
	}
}
