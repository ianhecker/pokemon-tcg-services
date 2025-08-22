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

func MakeCardIDs(IDs ...string) ([]CardID, error) {
	sanitized, err := SanitizeCardIDs(IDs...)
	if err != nil {
		return nil, err
	}
	var cardIDs []CardID
	for _, sani := range sanitized {
		cardIDs = append(cardIDs, CardID(sani))
	}
	return cardIDs, nil
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

func SanitizeCardIDs(IDs ...string) ([]string, error) {
	var sanitized []string
	var badIDs []string

	for _, ID := range IDs {
		sani, err := SanitizeCardID(ID)
		if err != nil {
			badIDs = append(badIDs, sani)
		} else {
			sanitized = append(sanitized, sani)
		}
	}
	if len(badIDs) > 0 {
		return nil, fmt.Errorf("invalid card IDs: %v", badIDs)
	}
	return sanitized, nil
}
