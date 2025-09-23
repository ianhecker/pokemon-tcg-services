package cards

import "github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"

type Printing string

type PrintingConditionPrices map[Printing]ConditionPrices

func MakePrintingConditionPrices(
	printings map[Printing]map[pokemontcg.Condition]Prices,
) PrintingConditionPrices {

	m := make(map[Printing]ConditionPrices)
	for printing, conditionMap := range printings {

		conditionPrices := MakeConditionPricesFromMap(conditionMap)
		m[printing] = conditionPrices
	}
	return m
}

func MakePrintingConditionPricesFromVariants(variants []Variant) PrintingConditionPrices {

	plinko := make(map[Printing]map[pokemontcg.Condition]Prices)
	for _, variant := range variants {

		printing := variant.Printing

		_, exists := plinko[printing]
		if !exists {
			plinko[printing] = make(map[pokemontcg.Condition]Prices)
		}

		condition := variant.Condition

		_, exists = plinko[printing][condition]
		if !exists {
			prices := MakePricesFromVariant(variant)
			plinko[printing][condition] = prices
		}
	}
	return MakePrintingConditionPrices(plinko)
}
