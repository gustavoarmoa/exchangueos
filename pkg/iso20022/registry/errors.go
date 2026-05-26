package registry

import "errors"

var (
	ErrDuplicate     = errors.New("registry: descriptor already registered")
	ErrNotFound      = errors.New("registry: descriptor not found")
	ErrEmptyField    = errors.New("registry: descriptor has empty required field")
	ErrMissingXSD    = errors.New("registry: descriptor missing XSDSourceURL")
	ErrNoRoute       = errors.New("router: no route matches request")
	ErrInvalidParty  = errors.New("router: invalid counterparty identifier")
	ErrUnsupportedOp = errors.New("router: operation not supported by target organisation")
)
