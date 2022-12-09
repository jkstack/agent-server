package scriptengine

import "errors"

// ErrTimeout timeout error
var ErrTimeout = errors.New("timeout")

// ErrUnexpectedError unexpected error
var ErrUnexpectedError = errors.New("unexpected")
