package pokemontcg

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Condition int

func (condition Condition) String() string {
	return ConditionToString(condition)
}

func (condition Condition) MarshalJSON() ([]byte, error) {
	return strconv.AppendQuote(nil, condition.String()), nil
}

func (condition *Condition) UnmarshalJSON(bytes []byte) error {
	var tmp string
	err := json.Unmarshal(bytes, &tmp)
	if err != nil {
		return fmt.Errorf("condition: error unmarshaling bytes: %w", err)
	}
	c, err := ParseCondition(tmp)
	if err != nil {
		return fmt.Errorf("condition: error parsing condition: %w", err)
	}
	*condition = c
	return nil
}

const (
	NearMint Condition = iota
	LightlyPlayed
	ModeratelyPlayed
	HeavilyPlayed
	Damaged
)

var AllConditions = []Condition{
	NearMint,
	LightlyPlayed,
	ModeratelyPlayed,
	HeavilyPlayed,
	Damaged,
}

func ConditionToString(condition Condition) string {
	switch condition {
	case NearMint:
		return "Near Mint"
	case LightlyPlayed:
		return "Lightly Played"
	case ModeratelyPlayed:
		return "Moderately Played"
	case HeavilyPlayed:
		return "Heavily Played"
	case Damaged:
		return "Damaged"
	default:
		return "Unknown"
	}
}

func ParseCondition(s string) (Condition, error) {
	switch s {
	case "Near Mint":
		return NearMint, nil
	case "Lightly Played":
		return LightlyPlayed, nil
	case "Moderately Played":
		return ModeratelyPlayed, nil
	case "Heavily Played":
		return HeavilyPlayed, nil
	case "Damaged":
		return Damaged, nil
	default:
		return Condition(-1), fmt.Errorf("unknown condition: '%s'", s)
	}
}
