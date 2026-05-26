package domain

import "errors"

var (
	ErrInvalidInput = errors.New("refdata: invalid input")
	ErrInactive     = errors.New("refdata: resource inactive")
	ErrExpired      = errors.New("refdata: ssi expired")
	ErrNotFound     = errors.New("refdata: not found")
	ErrStale        = errors.New("refdata: data stale")
)
