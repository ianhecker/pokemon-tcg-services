package pokemontcg

import (
	"fmt"
	"regexp"
	"strings"
)

var CardIDRegex = regexp.MustCompile(`^[A-Za-z0-9._:-]+$`)

type CardID string

func (cardID CardID) String() string {
	return string(cardID)
}

type Card struct {
	ID          CardID
	Name        string `json:"name"`
	Number      int    `json:"number,string"`
	Set         Set    `json:"set"`
	LastUpdated string `json:"lastUpdated"`
}

func MakeCardID(s string) (CardID, error) {
	ID := strings.TrimSpace(s)
	if ID == "" || !CardIDRegex.MatchString(ID) {
		return "", fmt.Errorf("invalid card ID: '%s'", ID)
	}
	return CardID(ID), nil
}
