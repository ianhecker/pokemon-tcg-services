package retry

import "fmt"

type RetryState int

func (state RetryState) String() string {
	return RetryStateToString(state)
}

const (
	Fail RetryState = iota
	Complete
	Yes
	Backoff
)

func RetryStateToString(state RetryState) string {
	switch state {
	case Fail:
		return "fail"
	case Complete:
		return "complete"
	case Yes:
		return "yes"
	case Backoff:
		return "backoff"
	default:
		return fmt.Sprintf("unknown state: %d", state)
	}
}
