package card

import (
	"errors"
	"fmt"
	"net/http"

	v2 "github.com/ianhecker/pokemon-tcg-services/internal/pokemontcgio/v2"
)

func ParseCardQuery(r *http.Request) ([]v2.CardID, error) {
	if r == nil {
		return nil, errors.New("request is nil")
	}

	q := r.URL.Query()
	queryStrings := q["id"]
	if len(queryStrings) == 0 {
		return nil, fmt.Errorf("missing required query parameter: id")
	}

	cardIDs, err := v2.MakeCardIDs(queryStrings...)
	if err != nil {
		return nil, fmt.Errorf("error sanitizing card IDs: %w", err)
	}

	if len(cardIDs) > 5 {
		return nil, fmt.Errorf("too many IDs (max 5)")
	}
	return cardIDs, nil
}
