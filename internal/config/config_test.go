package config_test

import (
	"testing"

	"github.com/ianhecker/pokemon-tcg-services/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestConfig_FlightCheck(t *testing.T) {
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

func TestLoad(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		token := "1234"
		t.Setenv(config.ENV_VAR_JUST_TCG_API_KEY, token)

		config, err := config.Load()
		require.NoError(t, err)

		gotToken := config.JustTCG.APIKey
		assert.Equal(t, token, gotToken.Reveal())
	})
	t.Run("missing Just TCG API key", func(t *testing.T) {
		_, err := config.Load()
		assert.ErrorContains(t, err, "Just TCG API Key is empty")
	})
}
