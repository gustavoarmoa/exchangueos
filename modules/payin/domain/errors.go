package domain

import "errors"

var (
	ErrInvalidInput      = errors.New("payin: invalid input")
	ErrInvalidTransition = errors.New("payin: invalid state transition")
	ErrDeadlineMissed    = errors.New("payin: deadline missed")
)
