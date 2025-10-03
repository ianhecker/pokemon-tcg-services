package proxy_test

import (
	"testing"
	"time"

	"github.com/ianhecker/pokemon-tcg-services/internal/networking/proxy"
	"github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {
	timer := proxy.Timer{}

	start := time.Now()
	timer.Start()
	timer.Stop()
	elapsed := time.Since(start)

	delta := elapsed - timer.Elapsed()

	assert.LessOrEqual(t, delta, 200*time.Millisecond)
	assert.GreaterOrEqual(t, delta, 10*time.Nanosecond)
}
