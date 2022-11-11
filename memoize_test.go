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

func TestSlice(t *testing.T) {
	called := 0
	work := func(x []int) []int {
		called++
		return x
	}

	m, err := Memo(work)
	assert.NoError(t, err)

	for i := 0; i < 1000; i++ {
		m([]int{1, 2, 3})
	}

	assert.Equal(t, 1, called)
}

func TestVariadic(t *testing.T) {
	called := 0
	work := func(x ...int) []int {
		called++
		return x
	}

	m, err := Memo(work)
	assert.NoError(t, err)

	for i := 0; i < 1000; i++ {
		switch i % 3 {
		case 0:
			m(0)
		case 1:
			m(0, 1)
		case 2:
			m(0, 1, 2)
		}
	}

	assert.Equal(t, 3, called)
}

func TestMatrix(t *testing.T) {
	called := 0
	work := func(x [4][4]int) [4][4]int {
		called++
		return x
	}

	m, err := Memo(work)
	assert.NoError(t, err)

	for i := 0; i < 1000; i++ {
		m([4][4]int{
			{1, 2, 3, 4},
			{5, 6, 7, 8},
			{9, 10, 11, 12},
			{13, 14, 15, 16},
		})
	}

	assert.Equal(t, 1, called)
}

func TestBadParamType(t *testing.T) {
	work := func(x map[string]int) map[string]int {
		return x
	}

	m, err := Memo(work)
	assert.NoError(t, err)

	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic")
			return
		}

		err := r.(error)

		assert.Error(t, err)
	}()

	_ = m(map[string]int{})

	t.Error("should have panicked")
}

func TestSliceVsArray(t *testing.T) {
	called := 0
	work := func(x interface{}) interface{} {
		called++
		return x
	}

	m, err := Memo(work)
	assert.NoError(t, err)

	m([]int{1, 2, 3})
	m([3]int{1, 2, 3})

	assert.Equal(t, 1, called)
}
