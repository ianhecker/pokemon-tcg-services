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

func MakeCardID(ID string) (CardID, error) {
	sani := strings.TrimSpace(ID)
	if sani == "" || !CardIDRegex.MatchString(sani) {
		return "", fmt.Errorf("invalid card ID: %s", sani)
	}
	return CardID(sani), nil
}
