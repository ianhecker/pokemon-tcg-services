package cards

import (
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

func MakeConditionPricesFromMap(m map[pokemontcg.Condition]Prices) ConditionPrices {
	var nm, lp, mp, hp, dmg Prices

	if prices, ok := m[pokemontcg.NearMint]; ok {
		nm = prices
	}
	if prices, ok := m[pokemontcg.LightlyPlayed]; ok {
		lp = prices
	}
	if prices, ok := m[pokemontcg.ModeratelyPlayed]; ok {
		mp = prices
	}
	if prices, ok := m[pokemontcg.HeavilyPlayed]; ok {
		hp = prices
	}
	if prices, ok := m[pokemontcg.Damaged]; ok {
		dmg = prices
	}
	return MakeConditionPrices(nm, lp, mp, hp, dmg)
}
