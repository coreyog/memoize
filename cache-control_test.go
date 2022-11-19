package memoize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheControlClear(t *testing.T) {
	called := 0
	work := func(x int) int {
		called++
		return x
	}

	m, cc, err := Memo(work)
	assert.NoError(t, err)

	m(1) // 1
	m(1)

	m(2) // 2
	m(2)

	cc.Clear()

	m(1) // 3
	m(1)

	m(2) // 4
	m(2)

	assert.Equal(t, 4, called)
}

func TestCacheControlRemove(t *testing.T) {
	called := 0
	work := func(x int, y string) int {
		called++
		return x
	}

	m, cc, err := Memo(work)
	assert.NoError(t, err)

	m(1, "x") // 1
	m(1, "x")

	m(1, "y") // 2
	m(1, "y")

	m(2, "x") // 3
	m(2, "x")

	cc.Remove(1)
	cc.Remove(3)

	m(1, "x") // 4
	m(1, "x")

	m(1, "y") // 5
	m(1, "y")

	m(2, "x")

	cc.Remove(1, "x")

	m(1, "x") // 6
	m(1, "x")

	m(1, "y")

	m(2, "x")

	cc.Remove("x")

	m(1, "x")
	m(1, "y")
	m(2, "x")

	cc.Remove(1, "x", 10.2) // does nothing, but shouldn't panic

	assert.Equal(t, 6, called)
}
