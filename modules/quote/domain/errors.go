package domain

import "errors"

var (
	ErrInvalidInput      = errors.New("quote: invalid input")
	ErrInvalidTransition = errors.New("quote: invalid state transition")
	ErrQuoteExpired      = errors.New("quote: quote expired")
	ErrQuoteRejected     = errors.New("quote: quote rejected")
)
