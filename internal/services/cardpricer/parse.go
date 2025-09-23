package cardpricer

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg/v1/cards"
)

func ParseCardQuery(r *http.Request) (cards.TCGPlayerID, error) {
	if r == nil {
		return "", errors.New("request is nil")
	}
	q := r.URL.Query()
	queryStrings := q["id"]
	if len(queryStrings) == 0 {
		return "", errors.New("missing required query: id")
	}
	ID, err := cards.MakeTCGPlayerID(queryStrings[0])
	if err != nil {
		return "", fmt.Errorf("error with query: %w", err)
	}
	return ID, nil
}
