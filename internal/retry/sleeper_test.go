package retry_test

import (
	"testing"
	"time"

	"github.com/ianhecker/pokemon-tcg-services/internal/retry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSleeper_Sleep(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		duration := 50 * time.Millisecond
		s := retry.NewSleeper(duration)

		start := time.Now()
		stop, done := s.Sleep()
		defer stop()

		<-done

		elapsed := time.Since(start)
		assert.GreaterOrEqual(t, elapsed, duration, "timer stopped too early")
	})

	t.Run("stop timer", func(t *testing.T) {
		d := 50 * time.Millisecond
		s := retry.NewSleeper(d)

		stop, done := s.Sleep()
		stopped := stop()
		require.True(t, stopped, "timer should have stopped")

		select {
		case <-done:
			assert.Fail(t, "timer channel fired after stop")
		case <-time.After(100 * time.Millisecond):
		}
	})
}
