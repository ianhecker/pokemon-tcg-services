package v2_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	v2 "github.com/ianhecker/pokemon-tcg-services/internal/pokemontcgio/v2"
)

func TestCardByID(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ID := "my-card-123"
		card, _ := v2.MakeCardID(ID)
		url := card.ToURL()
		assert.Equal(t, fmt.Sprintf("https://api.pokemontcg.io/v2/cards/%s", card), url)
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
			sanitized, err := v2.SanitizeCardID(test.ID)

			if test.valid {
				assert.Nil(t, err)
				assert.Equal(t, test.sanitized, sanitized)
			} else {
				assert.ErrorContains(t, err, "invalid card ID: "+test.ID)
			}
		})
	}
}
