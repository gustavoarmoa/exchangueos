package domain

import "errors"

var (
	ErrInvalidInput      = errors.New("compliance: invalid input")
	ErrInvalidTransition = errors.New("compliance: invalid state transition")
)
