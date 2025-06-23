package contravent

import "errors"

var ErrTimeout = errors.New("timeout")

type MatchError struct {
	Reasons []string
}

func (e MatchError) Error() string {
	return "event did not match schema"
}
