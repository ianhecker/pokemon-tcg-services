package retry

type RetryState int

const (
	No RetryState = iota
	Yes
	WithBackoff
)
