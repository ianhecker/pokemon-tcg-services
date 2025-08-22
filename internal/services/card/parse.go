package card

import (
	"errors"
	"fmt"
	"net/http"

	v2 "github.com/ianhecker/pokemon-tcg-services/internal/pokemontcgio/v2"
)

func ParseCardQuery(r *http.Request) (v2.CardID, error) {
	if r == nil {
		return "", errors.New("request is nil")
	}

	q := r.URL.Query()
	queryStrings := q["id"]
	if len(queryStrings) == 0 {
		return "", fmt.Errorf("missing required query parameter: id")
	}

	cardID, err := v2.MakeCardID(queryStrings[0])
	if err != nil {
		return "", fmt.Errorf("error sanitizing card ID: %w", err)
	}
	return cardID, nil
}
