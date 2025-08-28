package cards_test

import (
	"fmt"
	"testing"

	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg/v1/cards"
	"github.com/stretchr/testify/assert"
)

func TestCardID_MakeCardID(t *testing.T) {
	var tests = []struct {
		name      string
		ID        string
		sanitized string
		valid     bool
	}{
		{"happy path", "124014", "124014", true},
		{"happy path trim space", " 124014 ", "124014", true},
		{"empty", "", "", false},
		{"bad ID", "bad ID", "bad ID", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ID, err := cards.MakeCardID(test.ID)

			if test.valid {
				assert.Nil(t, err)
				assert.Equal(t, test.sanitized, ID.String())
			} else {
				assert.ErrorContains(t, err, fmt.Sprintf("invalid card ID: '%s'", test.ID))
			}
		})
	}
}

func TestCardID_UnmarshalJSON(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		pass  bool
	}{
		{"happy path", `"124014"`, true},
		{"empty", `""`, false},
		{"integer", `124014`, false},
		{"bad ID", `bad ID`, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			var cardID cards.CardID
			err := cardID.UnmarshalJSON([]byte(test.input))

			if test.pass {
				assert.NoError(t, err)
				assert.Equal(t, "124014", cardID.String())
			} else {
				assert.ErrorContains(t, err, "error unmarshaling card ID")
			}
		})
	}
}
