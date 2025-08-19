package retry

import (
	"time"
)

type RetryState int

const (
	NoRetry RetryState = iota
	RetryNoBackoff
	RetryWithBackoff
)

type RetryFunc func() bool

type Retryable struct {
	Attempt int
	State   RetryState
	Delay   time.Duration
	Do      RetryFunc
}

func (retry *Retryable) Retry() bool {
	retry.Attempt++
	time.Sleep(retry.Delay)

	return true
}
