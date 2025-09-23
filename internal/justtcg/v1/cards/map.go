package cards

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"
)

var ErrInvalid = errors.New("invalid")

func Validate(responseDTO ResponseDTO) error {
	if len(responseDTO.Data) == 0 {
		return errors.New("data field is empty")
	}

	card := responseDTO.Data[0]
	if card.TCGPlayerID == "" {
		return errors.New("card: TCGPlayerId is empty")
	}
	if len(card.Variants) == 0 {
		return errors.New("card: variants are empty")
	}

	for _, variant := range card.Variants {
		if variant.Printing == "" {
			return errors.New("card: variants: printings are empty")
		}
		if variant.Condition == "" {
			return errors.New("card: variants: condition is empty")
		}
	}
	return nil
}

func Map(responseDTO ResponseDTO) (Card, error) {
	err := Validate(responseDTO)
	if err != nil {
		return Card{}, fmt.Errorf("error validating response: %w", err)
	}

	cardDTO := responseDTO.Data[0]

	tcgPlayerID, err := MakeTCGPlayerID(string(cardDTO.TCGPlayerID))
	if err != nil {
		return Card{}, fmt.Errorf("error parsing TCGPlayerID: no regex match")
	}

	card := Card{
		ID:          cardDTO.ID,
		TCGPlayerID: tcgPlayerID,
		Name:        cardDTO.Name,
		Number:      cardDTO.Number,
		Rarity:      cardDTO.Rarity,
		Set:         cardDTO.Set,
	}

	lookupMap, err := VariantsToLookup(cardDTO.Variants)
	if err != nil {
		return Card{}, fmt.Errorf("error mapping variants: %w", err)
	}

	printings := make(map[Printing]Conditions)

	for printing, conditionsMap := range lookupMap {

		conditions := DenormalizeConditions(conditionsMap)
		printings[printing] = conditions
	}
	card.Prices = printings

	return card, nil
}

func MakeTCGPlayerID(ID string) (TCGPlayerID, error) {
	regex := func(s string) bool {
		s = strings.TrimSpace(s)
		return regexp.MustCompile(`^[0-9]+$`).MatchString(s)
	}
	match := regex(string(ID))
	if !match {
		return "", fmt.Errorf("error parsing TCGPlayerID: no regex match")
	}
	return TCGPlayerID(ID), nil
}

func MakePrice(priceDTO PriceDTO) *Price {
	if priceDTO == 0 {
		return nil
	}
	price := Price(priceDTO)
	return &price
}

func MapVariant(variant VariantDTO) (Printing, pokemontcg.Condition, *Prices, error) {
	printing := Printing(variant.Printing)

	condition, err := pokemontcg.ParseCondition(string(variant.Condition))
	if err != nil {
		return "", 0, nil, fmt.Errorf("error parsing variant condition: %w", err)
	}

	market := MakePrice(variant.Price)
	minPrice30d := MakePrice(variant.MinPrice30d)
	maxPrice30d := MakePrice(variant.MaxPrice30d)

	lastUpdated := time.Unix(variant.LastUpdated, 0).UTC()

	prices := Prices{
		Market:      market,
		MinPrice30d: minPrice30d,
		MaxPrice30d: maxPrice30d,
		LastUpdated: &lastUpdated,
	}
	return printing, condition, &prices, nil
}

func VariantsToLookup(variants []VariantDTO) (map[Printing]map[pokemontcg.Condition]*Prices, error) {

	doubleNestMap := make(map[Printing]map[pokemontcg.Condition]*Prices)

	doubleNest := func(printing Printing, condition pokemontcg.Condition, prices *Prices) {
		_, exists := doubleNestMap[printing]
		if !exists {
			doubleNestMap[printing] = make(map[pokemontcg.Condition]*Prices)
		}
		_, exists = doubleNestMap[printing][condition]
		if !exists {
			doubleNestMap[printing][condition] = prices
		}
	}

	for _, variant := range variants {
		printing, condition, prices, err := MapVariant(variant)
		if err != nil {
			return nil, fmt.Errorf("error mapping variant: %w", err)
		}
		doubleNest(printing, condition, prices)
	}
	return doubleNestMap, nil
}

func DenormalizeConditions(m map[pokemontcg.Condition]*Prices) Conditions {

	var conditions Conditions
	if prices, ok := m[pokemontcg.NearMint]; ok {
		conditions.NearMint = prices
	}
	if prices, ok := m[pokemontcg.LightlyPlayed]; ok {
		conditions.LightlyPlayed = prices
	}
	if prices, ok := m[pokemontcg.ModeratelyPlayed]; ok {
		conditions.ModeratelyPlayed = prices
	}
	if prices, ok := m[pokemontcg.HeavilyPlayed]; ok {
		conditions.HeavilyPlayed = prices
	}
	if prices, ok := m[pokemontcg.Damaged]; ok {
		conditions.Damaged = prices
	}
	return conditions
}
