package config

import "strings"

type Config struct {
}

type Token struct {
	token string
}

func NewToken(token string) *Token {
	return &Token{
		token: token,
	}
}

func (t Token) Inject() string {
	return string(t.token)
}

func (t Token) String() string {
	return t.obfuscate()
}

func (t Token) obfuscate() string {
	length := len(t.token)

	if length >= 10 {
		start := t.token[:3]
		middle := strings.Repeat("*", length-6)
		end := t.token[length-3:]

		return start + middle + end
	}
	return strings.Repeat("*", length)
}
