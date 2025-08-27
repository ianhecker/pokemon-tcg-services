package justtcg

import (
	"fmt"

	pokemon "github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"
)

func MakePricesAPIFromCardID(cardID pokemon.CardID) string {
	baseURL := "https://api.justtcg.com/v1"
	tcgplayerID := fmt.Sprintf("tcgplayerId=%s", cardID)
	return fmt.Sprintf("%s/cards?%s", baseURL, tcgplayerID)
}
