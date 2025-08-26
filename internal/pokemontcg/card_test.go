package pokemontcg_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"
)

func TestSanitizeCardID(t *testing.T) {
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
			ID, err := pokemontcg.MakeCardID(test.ID)

			if test.valid {
				assert.Nil(t, err)
				assert.Equal(t, test.sanitized, ID.String())
			} else {
				assert.ErrorContains(t, err, fmt.Sprintf("invalid card ID: '%s'", test.ID))
			}
		})
	}
}
