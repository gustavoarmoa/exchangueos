package domain

import "errors"

var (
	ErrInvalidInput = errors.New("risk: invalid input")
	ErrBreached     = errors.New("risk: limit breached")
)
