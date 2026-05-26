package domain

import "errors"

var (
	ErrInvalidInput         = errors.New("trade: invalid input")
	ErrInvalidTransition    = errors.New("trade: invalid state transition")
	ErrCancelReasonRequired = errors.New("trade: cancel reason required")
)
