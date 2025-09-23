package cards

import (
	"time"
)

// Data Transfer Objects (DTOs)
type ResponseDTO struct {
	Data []CardDTO `json:"data"`
}

type CardDTO struct {
	ID          string         `json:"id"`
	TCGPlayerID TCGPlayerIDDTO `json:"tcgplayerId"`
	Name        string         `json:"name"`
	Number      string         `json:"number"`
	Rarity      string         `json:"rarity"`
	Set         string         `json:"set"`
	Variants    []VariantDTO   `json:"variants"`
}

type TCGPlayerIDDTO string

type VariantDTO struct {
	Printing    PrintingDTO  `json:"printing"`
	Condition   ConditionDTO `json:"condition"`
	Price       PriceDTO     `json:"price"`
	MinPrice30d PriceDTO     `json:"minPrice30d"`
	MaxPrice30d PriceDTO     `json:"maxPrice30d"`
	LastUpdated int64        `json:"lastUpdated"`
}

type PriceDTO float64

type PrintingDTO string

type ConditionDTO string

// Domain Objects
type Cards []Card

type Card struct {
	ID          string      `json:"id"`
	TCGPlayerID TCGPlayerID `json:"tcgplayerId"`
	Name        string      `json:"name"`
	Number      string      `json:"number"`
	Rarity      string      `json:"rarity"`
	Set         string      `json:"set"`
	Prices      Printings   `json:"prices"`
}

type TCGPlayerID string

type Printing string

type Printings map[Printing]Conditions

type Conditions struct {
	NearMint         *Prices `json:"nearMint,omitempty"`
	LightlyPlayed    *Prices `json:"lightlyPlayed,omitempty"`
	ModeratelyPlayed *Prices `json:"moderatelyPlayed,omitempty"`
	HeavilyPlayed    *Prices `json:"heavilyPlayed,omitempty"`
	Damaged          *Prices `json:"damaged,omitempty"`
}

type Price float64

type Prices struct {
	Market      *Price     `json:"market,omitempty"`
	MinPrice30d *Price     `json:"minPrice30d,omitempty"`
	MaxPrice30d *Price     `json:"maxPrice30d,omitempty"`
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
}
