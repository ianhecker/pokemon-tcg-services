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

			token := config.MakeToken(test.input)
			assert.Equal(t, test.output, token.String())
		})
	}
}

func TestToken_Reveal(t *testing.T) {
	var tests = []struct {
		name   string
		secret string
	}{
		{"empty token", ""},
		{"short token", "123"},
		{"long token", "abcdefghijklmnopqrstuvwxyz"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			token := config.MakeToken(test.secret)
			assert.Equal(t, test.secret, token.Reveal())
		})
	}
}

func TestConfg_FlightCheck(t *testing.T) {
	var tests = []struct {
		name     string
		token    string
		hasError bool
	}{
		{"happy path", "1234", false},
		{"empty token", "", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			config := config.Config{
				JustTCG: struct{ APIKey config.Token }{
					APIKey: config.MakeToken(test.token),
				},
			}
			err := config.FlightCheck()
			assert.Equal(t, test.hasError, err != nil)
		})
	}
}
