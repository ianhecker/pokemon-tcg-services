package pokemontcg_test

import (
	"encoding/json"
	"testing"

	"github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCondition_ParseCondition(t *testing.T) {
	var tests = []struct {
		input     string
		condition pokemontcg.Condition
	}{
		{"Near Mint", pokemontcg.NearMint},
		{"Lightly Played", pokemontcg.LightlyPlayed},
		{"Moderately Played", pokemontcg.ModeratelyPlayed},
		{"Heavily Played", pokemontcg.HeavilyPlayed},
		{"Damaged", pokemontcg.Damaged},
		{"Unknown", -1},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {

			condition, err := pokemontcg.ParseCondition(test.input)
			if test.input != "Unknown" {
				assert.NoError(t, err)
				assert.Equal(t, test.condition, condition)
			} else {
				assert.ErrorContains(t, err, "unknown condition")
			}
		})
	}
}

func TestCondition_MarshalJSON(t *testing.T) {
	var tests = []struct {
		condition pokemontcg.Condition
		expected  string
	}{
		{pokemontcg.NearMint, `"Near Mint"`},
		{pokemontcg.LightlyPlayed, `"Lightly Played"`},
		{pokemontcg.ModeratelyPlayed, `"Moderately Played"`},
		{pokemontcg.HeavilyPlayed, `"Heavily Played"`},
		{pokemontcg.Damaged, `"Damaged"`},
		{-1, `"Unknown"`},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {

			bytes, err := test.condition.MarshalJSON()
			if test.expected != "Unknown" {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, string(bytes))
			} else {
				assert.ErrorContains(t, err, "error marshaling condition")
			}
		})
	}
}

func TestCondition_UnmarshalJSON(t *testing.T) {
	var tests = []struct {
		input    string
		expected pokemontcg.Condition
	}{
		{"Near Mint", pokemontcg.NearMint},
		{"Lightly Played", pokemontcg.LightlyPlayed},
		{"Moderately Played", pokemontcg.ModeratelyPlayed},
		{"Heavily Played", pokemontcg.HeavilyPlayed},
		{"Damaged", pokemontcg.Damaged},
		{"Unknown", -1},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {

			bytes, err := json.Marshal(test.input)
			require.NoError(t, err)

			var condition pokemontcg.Condition
			err = condition.UnmarshalJSON(bytes)

			if test.input != "Unknown" {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, condition)
			} else {
				assert.ErrorContains(t, err, "error unmarshaling condition")
			}
		})
	}
}
