package domain

import "errors"

var (
	ErrInvalidInput      = errors.New("admin: invalid input")
	ErrInvalidTransition = errors.New("admin: invalid state transition")
)
