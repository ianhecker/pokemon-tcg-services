package config_test

import (
	"testing"

	"github.com/ianhecker/pokemon-tcg-services/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestToken_String(t *testing.T) {
	var tests = []struct {
		name   string
		input  string
		output string
	}{
		{"vary short token", "abc", "***"},
		{"short token", "abcdefghij", "abc****hij"},
		{"long token", "abcdefghijklmnopqrstuvwxyz", "abc********************xyz"},
		{"empty", "", ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			token := config.NewToken(test.input)
			assert.Equal(t, test.output, token.String())
		})
	}
}
