package domain

import "errors"

var (
	ErrInvalidInput      = errors.New("cfets_confirmation: invalid input")
	ErrInvalidTransition = errors.New("cfets_confirmation: invalid state transition")
)
