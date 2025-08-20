package retry

import "time"

type Sleepable interface {
	Sleep() (stop func() bool, done <-chan time.Time)
	Duration() time.Duration
}

type Sleeper time.Duration

func NewSleeper(duration time.Duration) Sleepable {
	s := Sleeper(duration)
	return &s
}

func (s *Sleeper) Duration() time.Duration {
	return time.Duration(*s)
}

func (s *Sleeper) Sleep() (func() bool, <-chan time.Time) {
	timer := time.NewTimer(s.Duration())
	return timer.Stop, timer.C
}
