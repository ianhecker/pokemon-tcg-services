package proxy

import "time"

type Timer struct {
	start   time.Time
	elapsed time.Duration
}

func (t *Timer) Start() {
	t.start = time.Now()
}

func (t *Timer) Stop() {
	t.elapsed = time.Since(t.start)
}

func (t *Timer) Elapsed() time.Duration {
	return t.elapsed
}
