package domain

import "errors"

var (
	ErrInvalidInput      = errors.New("cls_settlement: invalid input")
	ErrInvalidTransition = errors.New("cls_settlement: invalid state transition")
	ErrCycleClosed       = errors.New("cls_settlement: cycle closed")
)
