package domain

import "errors"

var (
	ErrInvalidInput      = errors.New("cfets_capture: invalid input")
	ErrInvalidTransition = errors.New("cfets_capture: invalid state transition")
)
