package cards

import (
	"errors"

	"github.com/ianhecker/pokemon-tcg-services/internal/pokemontcg"
)

type ConditionPrices struct {
	NearMint         Prices `json:"nearMint,omitempty"`
	LightlyPlayed    Prices `json:"lightlyPlayed,omitempty"`
	ModeratelyPlayed Prices `json:"moderatelyPlayed,omitempty"`
	HeavilyPlayed    Prices `json:"heavilyPlayed,omitempty"`
	Damaged          Prices `json:"damaged,omitempty"`
}

func MakeConditionPrices(nm, lp, mp, hp, dmg Prices) ConditionPrices {
	return ConditionPrices{
		NearMint:         nm,
		LightlyPlayed:    lp,
		ModeratelyPlayed: mp,
		HeavilyPlayed:    hp,
		Damaged:          dmg,
	}
}

func MakeConditionPricesFromVariants(
	m map[pokemontcg.Condition]Prices,
) (ConditionPrices, error) {
	var nm, lp, mp, hp, dmg Prices
	counter := 0

	if prices, ok := m[pokemontcg.NearMint]; ok {
		nm = prices
		counter++
	}
	if prices, ok := m[pokemontcg.LightlyPlayed]; ok {
		lp = prices
		counter++
	}
	if prices, ok := m[pokemontcg.ModeratelyPlayed]; ok {
		mp = prices
		counter++
	}
	if prices, ok := m[pokemontcg.HeavilyPlayed]; ok {
		hp = prices
		counter++
	}
	if prices, ok := m[pokemontcg.Damaged]; ok {
		dmg = prices
		counter++
	}
	if counter == 0 {
		return ConditionPrices{}, errors.New("given zero prices")
	}

	return MakeConditionPrices(nm, lp, mp, hp, dmg), nil
}
