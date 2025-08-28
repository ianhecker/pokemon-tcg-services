package cards

import (
	"encoding/json"
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

func (cardID *CardID) UnmarshalJSON(bytes []byte) error {
	var tmp string
	err := json.Unmarshal(bytes, &tmp)
	if err != nil {
		return fmt.Errorf("error unmarshaling card ID: %w", err)
	}
	ID, err := MakeCardID(tmp)
	if err != nil {
		return fmt.Errorf("error unmarshaling card ID: %w", err)
	}
	*cardID = ID
	return nil
}
