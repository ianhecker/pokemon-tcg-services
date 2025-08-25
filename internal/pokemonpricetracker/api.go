package pokemonpricetracker

import (
	"fmt"

	"github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"
)

func MakePricesAPIFromCardID(ID pokemontcg.CardID) string {
	return fmt.Sprintf("https://www.pokemonpricetracker.com/api/prices?id=%s", ID.String())
}
