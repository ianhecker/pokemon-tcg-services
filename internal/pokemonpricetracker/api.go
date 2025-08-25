package pokemonpricetracker

import (
	"fmt"

	"github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"
)

func MakePricesAPI(card pokemontcg.Card) string {
	return fmt.Sprintf("https://www.pokemonpricetracker.com/api/prices?id=%s", card.ID)
}
