package config

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

type Token struct {
	token string
}

func MakeToken(token string) Token {
	return Token{
		token: token,
	}
}

func (t Token) Reveal() string {
	return t.token
}

func (t Token) String() string {
	return t.obfuscate()
}

func (t Token) obfuscate() string {
	length := len(t.token)
	if length >= 10 {
		return t.token[:3] + strings.Repeat("*", length-6) + t.token[length-3:]
	}
	return strings.Repeat("*", length)
}

type Config struct {
	JustTCG struct {
		APIKey Token
	}
}

func (config Config) FlightCheck() error {
	if config.JustTCG.APIKey.Reveal() == "" {
		return errors.New("config flight check: Just TCG API Key is empty")
	}
	return nil
}

func Load() (Config, error) {
	_ = viper.BindEnv("JUST_TCG_API_KEY")
	viper.AutomaticEnv()

	var config Config
	token := MakeToken(viper.GetString("JUST_TCG_API_KEY"))
	config.JustTCG.APIKey = token

	return config, nil
}
