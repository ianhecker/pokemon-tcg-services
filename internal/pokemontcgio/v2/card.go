package v2

import (
	"fmt"
	"regexp"
	"strings"
)

const CardByIDHttpRequestFmt string = "https://api.pokemontcg.io/v2/cards/%s"

var CardIDRegex = regexp.MustCompile(`^[A-Za-z0-9._:-]+$`)

type CardID string

func MakeCardID(ID string) (CardID, error) {
	cardID, err := SanitizeCardID(ID)
	if err != nil {
		return "", err
	}
	return CardID(cardID), nil
}

func (card CardID) ToURL() string {
	return fmt.Sprintf(CardByIDHttpRequestFmt, string(card))
}

func SanitizeCardID(ID string) (string, error) {
	sani := strings.TrimSpace(ID)
	if sani == "" || !CardIDRegex.MatchString(sani) {
		return "", fmt.Errorf("invalid card ID: %s", sani)
	}
	return sani, nil
}
