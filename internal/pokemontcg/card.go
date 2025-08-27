package pokemontcg

import (
	"fmt"
	"regexp"
	"strings"
)

var CardIDRegex = regexp.MustCompile(`^[0-9]+$`)

type CardID string

func (cardID CardID) String() string {
	return string(cardID)
}

type Card struct {
	ID          string             `json:"id"`
	TCGPlayerID CardID             `json:"tcgplayerId"`
	Name        string             `json:"name"`
	Number      string             `json:"number"`
	Rarity      string             `json:"rarity"`
	Set         string             `json:"set"`
	Pricing     []ConditionPricing `json:"variants"`
}

func MakeCardID(s string) (CardID, error) {
	ID := strings.TrimSpace(s)
	if ID == "" || !CardIDRegex.MatchString(ID) {
		return "", fmt.Errorf("invalid card ID: '%s'", ID)
	}
	return CardID(ID), nil
}
