package cards

import (
	"fmt"
)

const baseURL = "https://api.justtcg.com/v1"

func GetCardByID(ID CardID) string {
	return fmt.Sprintf("%s/cards?tcgplayerId=%s", baseURL, ID)
}
