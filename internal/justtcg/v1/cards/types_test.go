package cards_test

import (
	"time"

	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg/v1/cards"
)

var expected = struct {
	Response cards.ResponseDTO
	Card     cards.Card
}{
	Response: _responseDTO,
	Card:     _card,
}

var _responseDTO = cards.ResponseDTO{
	Data: []cards.CardDTO{
		{
			ID:          "pokemon-base-set-shadowless-charizard-holo-rare",
			TCGPlayerID: cards.TCGPlayerIDDTO("106999"),
			Name:        "Charizard",
			Number:      "004/102",
			Rarity:      "Holo Rare",
			Set:         "Base Set (Shadowless)",
			Variants: []cards.VariantDTO{
				{
					Printing:    cards.PrintingDTO("1st Edition Holofoil"),
					Condition:   cards.ConditionDTO("Heavily Played"),
					Price:       cards.PriceDTO(225),
					LastUpdated: 1743656485,
				},
				{
					Printing:    cards.PrintingDTO("1st Edition Holofoil"),
					Condition:   cards.ConditionDTO("Near Mint"),
					Price:       cards.PriceDTO(503.49),
					MinPrice30d: cards.PriceDTO(500.99),
					MaxPrice30d: cards.PriceDTO(503.49),
					LastUpdated: 1758706033,
				},
				{
					Printing:    cards.PrintingDTO("Unlimited Holofoil"),
					Condition:   cards.ConditionDTO("Damaged"),
					Price:       cards.PriceDTO(524.96),
					MinPrice30d: cards.PriceDTO(523.77),
					MaxPrice30d: cards.PriceDTO(533.29),
					LastUpdated: 1758706033,
				},
				{
					Printing:    cards.PrintingDTO("Unlimited Holofoil"),
					Condition:   cards.ConditionDTO("Heavily Played"),
					Price:       cards.PriceDTO(645.05),
					MinPrice30d: cards.PriceDTO(642.95),
					MaxPrice30d: cards.PriceDTO(649),
					LastUpdated: 1758706033,
				},
				{
					Printing:    cards.PrintingDTO("Unlimited Holofoil"),
					Condition:   cards.ConditionDTO("Moderately Played"),
					Price:       cards.PriceDTO(764.78),
					MinPrice30d: cards.PriceDTO(761.02),
					MaxPrice30d: cards.PriceDTO(813.73),
					LastUpdated: 1758706033,
				},
				{
					Printing:    cards.PrintingDTO("Unlimited Holofoil"),
					Condition:   cards.ConditionDTO("Lightly Played"),
					Price:       cards.PriceDTO(1299),
					MinPrice30d: cards.PriceDTO(1299),
					MaxPrice30d: cards.PriceDTO(1299),
					LastUpdated: 1758706033,
				},
				{
					Printing:    cards.PrintingDTO("Unlimited Holofoil"),
					Condition:   cards.ConditionDTO("Near Mint"),
					Price:       cards.PriceDTO(1923.19),
					MinPrice30d: cards.PriceDTO(1700),
					MaxPrice30d: cards.PriceDTO(1923.19),
					LastUpdated: 1758706033,
				},
				{
					Printing:    cards.PrintingDTO("1st Edition Holofoil"),
					Condition:   cards.ConditionDTO("Moderately Played"),
					Price:       cards.PriceDTO(5495),
					MinPrice30d: cards.PriceDTO(5495),
					MaxPrice30d: cards.PriceDTO(5495),
					LastUpdated: 1754787579,
				},
			},
		},
	},
}

var _card = cards.Card{
	ID:          "pokemon-base-set-shadowless-charizard-holo-rare",
	TCGPlayerID: "106999",
	Name:        "Charizard",
	Number:      "004/102",
	Rarity:      "Holo Rare",
	Set:         "Base Set (Shadowless)",
	Prices: cards.Printings{
		"1st Edition Holofoil": cards.Conditions{
			NearMint: &cards.Prices{
				Market:      fptr(503.49),
				MinPrice30d: fptr(500.99),
				MaxPrice30d: fptr(503.49),
				LastUpdated: tptr(time.Unix(1758706033, 0).UTC()),
			},
			ModeratelyPlayed: &cards.Prices{
				Market:      fptr(5495),
				MinPrice30d: fptr(5495),
				MaxPrice30d: fptr(5495),
				LastUpdated: tptr(time.Unix(1754787579, 0).UTC()),
			},
			HeavilyPlayed: &cards.Prices{
				Market:      cards.MakePrice(225),
				LastUpdated: tptr(time.Unix(1743656485, 0).UTC()),
			},
		},
		"Unlimited Holofoil": cards.Conditions{
			NearMint: &cards.Prices{
				Market:      fptr(1923.19),
				MinPrice30d: fptr(1700),
				MaxPrice30d: fptr(1923.19),
				LastUpdated: tptr(time.Unix(1758706033, 0).UTC()),
			},
			LightlyPlayed: &cards.Prices{
				Market:      fptr(1299),
				MinPrice30d: fptr(1299),
				MaxPrice30d: fptr(1299),
				LastUpdated: tptr(time.Unix(1758706033, 0).UTC()),
			},
			ModeratelyPlayed: &cards.Prices{
				Market:      fptr(764.78),
				MinPrice30d: fptr(761.02),
				MaxPrice30d: fptr(813.73),
				LastUpdated: tptr(time.Unix(1758706033, 0).UTC()),
			},
			HeavilyPlayed: &cards.Prices{
				Market:      fptr(645.05),
				MinPrice30d: fptr(642.95),
				MaxPrice30d: fptr(649),
				LastUpdated: tptr(time.Unix(1758706033, 0).UTC()),
			},
			Damaged: &cards.Prices{
				Market:      fptr(524.96),
				MinPrice30d: fptr(523.77),
				MaxPrice30d: fptr(533.29),
				LastUpdated: tptr(time.Unix(1758706033, 0).UTC()),
			},
		},
	},
}
