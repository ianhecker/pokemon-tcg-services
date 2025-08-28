package cards

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

func MakeCardID(s string) (CardID, error) {
	ID := strings.TrimSpace(s)
	if ID == "" || !CardIDRegex.MatchString(ID) {
		return "", fmt.Errorf("invalid card ID: '%s'", ID)
	}
	return CardID(ID), nil
}
