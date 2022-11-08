package memoize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadFuncErrors(t *testing.T) {
	_, err := Memo("not a func")
	assert.Equal(t, err, ErrNotAFunc)

	noArgs := func() bool { return true }
	_, err = Memo(noArgs)
	assert.Equal(t, err, ErrMissingArgs)

	noReturns := func(x int) {}
	_, err = Memo(noReturns)
	assert.Equal(t, err, ErrMissingReturns)
}

func TestMultiCall(t *testing.T) {
	called := 0
	square := func(x int) int {
		called++
		return x * x
	}

	m, err := Memo(square)
	assert.NoError(t, err)

	for i := 0; i < 100; i++ {
		mResult := m(i)
		m2Result := m(i)
		squareResult := square(i)

		assert.Equal(t, mResult, m2Result)
		assert.Equal(t, mResult, squareResult)
	}

	assert.Equal(t, called, 200)
}

func TestSimple(t *testing.T) {
	called := 0
	work := func(x int) int {
		called++
		return x
	}

	m, err := Memo(work)
	assert.NoError(t, err)

	for i := 0; i < 1000; i++ {
		m(0)
	}

	assert.Equal(t, 1, called)
}

func TestManyArgsManyRets(t *testing.T) {
	called := 0
	work := func(x int, y string, z float64) (int, string, float64) {
		called++
		return x, y, z
	}

	m, err := Memo(work)
	assert.NoError(t, err)

	for i := 0; i < 1000; i++ {
		m(0, "x", 3.14)
	}

	assert.Equal(t, 1, called)
}
