package memoize

import "github.com/pkg/errors"

var (
	ErrNotAFunc       = errors.New("not a function")
	ErrMissingArgs    = errors.New("target function must accept at least 1 argument")
	ErrMissingReturns = errors.New("target function must return at least 1 value")
)
