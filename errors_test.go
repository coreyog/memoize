package memoize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadFuncErrors(t *testing.T) {
	_, _, err := Memo("not a func")
	assert.Equal(t, err, ErrNotAFunc)

	noArgs := func() bool { return true }
	_, _, err = Memo(noArgs)
	assert.Equal(t, err, ErrMissingArgs)

	noReturns := func(x int) {}
	_, _, err = Memo(noReturns)
	assert.Equal(t, err, ErrMissingReturns)
}
