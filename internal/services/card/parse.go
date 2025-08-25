package card

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"
)

func ParseCardQuery(r *http.Request) (pokemontcg.Card, error) {
	if r == nil {
		return pokemontcg.Card{}, errors.New("request is nil")
	}

	q := r.URL.Query()
	queryStrings := q["id"]
	if len(queryStrings) == 0 {
		return pokemontcg.Card{}, fmt.Errorf("missing required query parameter: id")
	}

	card, err := pokemontcg.MakeCard(queryStrings[0])
	if err != nil {
		return pokemontcg.Card{}, fmt.Errorf("error sanitizing card ID: %w", err)
	}
	return card, nil
}
