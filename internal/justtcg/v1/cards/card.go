package cards

import (
	"encoding/json"
	"fmt"
)

type Card struct {
	ID          string                  `json:"id"`
	TCGPlayerID CardID                  `json:"tcgplayerId"`
	Name        string                  `json:"name"`
	Number      string                  `json:"number"`
	Rarity      string                  `json:"rarity"`
	Set         string                  `json:"set"`
	Prices      PrintingConditionPrices `json:"prices"`
}

func MakeCard(
	ID string,
	tcgPlayerID CardID,
	name string,
	number string,
	rarity string,
	set string,
	prices PrintingConditionPrices,
) Card {
	return Card{
		ID:          ID,
		TCGPlayerID: tcgPlayerID,
		Name:        name,
		Number:      number,
		Rarity:      rarity,
		Set:         set,
		Prices:      prices,
	}
}

type RawCard struct {
	ID          string    `json:"id"`
	TCGPlayerID CardID    `json:"tcgplayerId"`
	Name        string    `json:"name"`
	Number      string    `json:"number"`
	Rarity      string    `json:"rarity"`
	Set         string    `json:"set"`
	Variants    []Variant `json:"variants"`
}

func (card *Card) UnmarshalJSON(bytes []byte) error {
	var raw RawCard
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("error unmarshaling card: %w", err)
	}

	prices := MakePrintingConditionPricesFromVariants(raw.Variants)

	*card = MakeCard(
		raw.ID,
		raw.TCGPlayerID,
		raw.Name,
		raw.Number,
		raw.Rarity,
		raw.Set,
		prices,
	)
	return nil
}
