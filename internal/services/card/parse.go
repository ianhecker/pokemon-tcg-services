package card

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"
)

func ParseCardQuery(r *http.Request) (pokemontcg.CardID, error) {
	if r == nil {
		return "", errors.New("request is nil")
	}

	q := r.URL.Query()
	queryStrings := q["id"]
	if len(queryStrings) == 0 {
		return "", fmt.Errorf("missing required query parameter: id")
	}

	ID, err := pokemontcg.MakeCardID(queryStrings[0])
	if err != nil {
		return "", fmt.Errorf("error sanitizing card ID: %w", err)
	}
	return ID, nil
}
