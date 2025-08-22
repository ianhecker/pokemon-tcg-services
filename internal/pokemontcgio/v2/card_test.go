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
			} else {
				assert.ErrorContains(t, err, "invalid card ID: "+test.ID)
			}
			assert.Equal(t, test.sanitized, sanitized)
		})
	}
}

func TestCardByID_AreValidCardIDs(t *testing.T) {
	var tests = []struct {
		name     string
		IDs      []string
		expected []string
		valid    bool
	}{
		{"happy path", []string{"xy1-1"}, []string{"xy1-1"}, true},
		{"happy path list", []string{"xy1-1", "xy2-2"}, []string{"xy1-1", "xy2-2"}, true},
		{"happy path trim space", []string{" xy1-1 "}, []string{"xy1-1"}, true},
		{"happy path trim space list", []string{" xy1-1 ", " xy2-2 "}, []string{"xy1-1", "xy2-2"}, true},
		{"empty", []string{""}, []string{""}, false},
		{"bad ID", []string{"bad ID"}, []string{"bad ID"}, false},
		{"bad ID list", []string{"bad ID 1", "bad ID 2"}, []string{"bad ID 1", "bad ID 2"}, false},
		{"bad ID w valid", []string{"base1-1", "bad ID 1", "bad ID 2"}, []string{"bad ID 1", "bad ID 2"}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			IDs, err := v2.SanitizeCardIDs(test.IDs...)

			if test.valid {
				assert.Nil(t, err)
				assert.Equal(t, test.expected, IDs)
			} else {
				assert.ErrorContains(t, err, fmt.Sprintf("invalid card IDs: %v", IDs))
				assert.Equal(t, test.expected, IDs)
			}
		})
	}
}
