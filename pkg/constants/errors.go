package constants

import "errors"

const (
	ErrMsgFailedSerializeRecord = "failed to serialized record"
	ErrMsgEmptyId               = "id could not be empty"
)

var (
	ErrPanic = errors.New("panic occurs")
)
