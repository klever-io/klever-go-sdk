package builders

import "errors"

// ErrInvalidValue signals that an invalid value was provided
var ErrInvalidValue = errors.New("invalid value")

// ErrNilValue signals that a nil value was provided
var ErrNilValue = errors.New("nil value")

// ErrNilAddress signals that a nil address was provided
var ErrNilAddress = errors.New("nil address handler")

// ErrInvalidAddress signals that an invalid address was provided
var ErrInvalidAddress = errors.New("invalid address handler")
