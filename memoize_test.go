package memoize

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Model struct {
	ID     int
	Called *int
}

func NewModel() *Model {
	return &Model{ID: 1, Called: new(int)}
}

func (m *Model) PointerReceiver(x int) int {
	*m.Called++
	return m.ID
}

func (m Model) NonPointerReceiver(x int) int {
	*m.Called++
	return m.ID
}

func TestSimple(t *testing.T) {
	called := 0
	work := func(x int) int {
		called++
		return x
	}

	m, _, err := Memo(work)
	assert.NoError(t, err)

	for i := 0; i < 1000; i++ {
		m(0)
	}

	assert.Equal(t, 1, called)
}

func TestMultiCall(t *testing.T) {
	called := 0
	square := func(x int) int {
		called++
		return x * x
	}

	m, _, err := Memo(square)
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

func TestManyArgsManyRets(t *testing.T) {
	called := 0
	work := func(x int, y string, z float64) (int, string, float64) {
		called++
		return x, y, z
	}

	m, _, err := Memo(work)
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

	m, _, err := Memo(work)
	assert.NoError(t, err)

	for i := 0; i < 1000; i++ {
		m([]int{1, 2, 3})
	}

	assert.Equal(t, 1, called)
}

func TestNilSlice(t *testing.T) {
	called := 0
	work := func(x []int) []int {
		called++
		return x
	}

	m, _, err := Memo(work)
	assert.NoError(t, err)

	m(nil) // 1
	m(nil)

	assert.Equal(t, 1, called)
}

func TestVariadic(t *testing.T) {
	called := 0
	work := func(x ...int) []int {
		called++
		return x
	}

	m, _, err := Memo(work)
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

	m, _, err := Memo(work)
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

	m, _, err := Memo(work)
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

	m, _, err := Memo(work)
	assert.NoError(t, err)

	m([]int{1, 2, 3})
	m([3]int{1, 2, 3})

	assert.Equal(t, 1, called)
}

func TestMapReturn(t *testing.T) {
	called := 0
	work := func(x int) map[int]int {
		called++
		return map[int]int{x: x}
	}

	m, _, err := Memo(work)
	assert.NoError(t, err)

	map1 := m(1)
	map2 := m(1)

	if !reflect.DeepEqual(map1, map2) {
		t.Error("results not equal")
	}

	assert.Equal(t, 1, called)
}

func TestRecursion(t *testing.T) {
	called := 0

	var f func(int) int   // predefine
	f = func(x int) int { // initialize
		called++
		if x <= 1 {
			return 1
		}

		return x * f(x-1) // factorial
	}

	f, _, err := Memo(f) // overwrite
	assert.NoError(t, err)

	x := f(10) // calls f(10), f(9), f(8), ..., f(1)
	assert.Equal(t, 3628800, x)
	assert.Equal(t, 10, called)

	x = f(11) // only calls f(11)
	assert.Equal(t, 39916800, x)
	assert.Equal(t, 11, called)

	x = f(5) // no new calls
	assert.Equal(t, 120, x)
	assert.Equal(t, 11, called)
}

func TestMemberFuncs(t *testing.T) {
	model := NewModel()
	assert.NotNil(t, model.Called)
	assert.Equal(t, 0, *model.Called)

	m1, _, err := Memo(model.PointerReceiver)
	assert.NoError(t, err)

	m2, _, err := Memo(model.NonPointerReceiver)
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		m1(0)
		m2(0)
	}

	assert.Equal(t, 2, *model.Called)
}

func TestNonErrorPanic(t *testing.T) {
	work := func(x int) int {
		panic("test panic")
	}

	m, _, err := Memo(work)
	assert.NoError(t, err)

	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic")
			return
		}

		err, _ := r.(error)

		assert.NoError(t, err)
		assert.Equal(t, "test panic", r)
	}()

	m(0)
}
