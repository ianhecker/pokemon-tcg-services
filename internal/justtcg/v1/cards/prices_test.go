package cards_test

import (
	"testing"
	"time"

	"github.com/ianhecker/pokemon-tcg-services/internal/justtcg/v1/cards"
	"github.com/stretchr/testify/assert"
)

func TestPrices_MakePrices(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		p := cards.MakePrices(
			1.1,
			1.2,
			1.3,
			1756329181,
		)
		assert.Equal(t, time.Unix(1756329181, 0).UTC(), p.LastUpdated)
	})
}
