package pokemontcg

import (
	"fmt"
	"regexp"
	"strings"
)

var CardIDRegex = regexp.MustCompile(`^[A-Za-z0-9._:-]+$`)

type Card struct {
	ID     string
	Name   string `json:"name"`
	Number int    `json:"number"`
	Set    Set    `json:"set"`
}

func MakeCard(ID string) (Card, error) {
	cardID, err := SanitizeID(ID)
	if err != nil {
		return Card{}, err
	}
	return Card{
		ID: cardID,
	}, nil
}

func SanitizeID(ID string) (string, error) {
	sani := strings.TrimSpace(ID)
	if sani == "" || !CardIDRegex.MatchString(sani) {
		return "", fmt.Errorf("invalid card ID: %s", sani)
	}
	return sani, nil
}
