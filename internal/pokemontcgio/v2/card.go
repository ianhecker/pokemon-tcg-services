package v2

import (
	"fmt"
	"regexp"
	"strings"
)

const CardByIDFmt string = "api.pokemontcg.io/v2/cards/%s"

var CardIDRegex = regexp.MustCompile(`^[A-Za-z0-9._:-]+$`)

func CardByID(ID string) string {
	return fmt.Sprintf(CardByIDFmt, ID)
}

func SanitizeCardID(ID string) (string, error) {
	sani := strings.TrimSpace(ID)
	if sani == "" || !CardIDRegex.MatchString(sani) {
		return sani, fmt.Errorf("invalid card ID: %s", sani)
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
		return badIDs, fmt.Errorf("invalid card IDs: %v", badIDs)
	}
	return sanitized, nil
}
