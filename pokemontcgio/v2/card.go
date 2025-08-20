package v2

import (
	"fmt"
)

const CardByIDFmt string = "api.pokemontcg.io/v2/cards/%s"

func CardByID(ID string) string {
	return fmt.Sprintf(CardByIDFmt, ID)
}
