# Memoize

A generic memoization library.

Needs more testing before going to prod.

## Example

```
func IsPrime(val int) bool {
	val = int(math.Abs(float64(val)))
	if (val&1) == 0 || val == 1 {
		return val == 2
	}

	sqrt := int(math.Ceil(math.Sqrt(float64(val))))
	for i := 3; i <= sqrt; i += 2 {
		if (val % i) == 0 {
			return false
		}
	}

	return true
}

...

mIsPrime, _ := Memoize.Memo(IsPrime)

// no change in signature
// mIsPrime = func(int) bool

IsPrime(2) // called, returned true
IsPrime(4) // called, returned false
IsPrime(2) // called, returned true
IsPrime(4) // called, returned false

mIsPrime(2) // called, returned true
mIsPrime(4) // called, returned false
mIsPrime(2) // not called, returned true
mIsPrime(4) // not called, returned false
```

Not every function can be memoized. The types of allowed parameters are limited
by this library. Returned types don't have any restrictions.

### Do's:
* primitives
* 1D slices
* ND arrays
* Variadic
* structs containing only Do's

### Don'ts:
* maps
* channels
* functions
* ND slices
* pointers

## Notes
To support 1D slices, the slices are converted to Arrays before being placed in
the cache. This results in the following edge case:

```
called := 0
work := func(x interface{}) interface{} {
  called++
  return x
}

m, _ := Memo(work)

m([]int{1, 2, 3}) // no record of params, work runs
m([3]int{1, 2, 3}) // "collides" with first call, work won't run

// called = 1
```

Because I want to return a func with the exact same parameter list and return
values as the provided input, any errors that arise will panic.

Member functions work but are probably a bad idea. Oh well, I won't stand in the
way.

I used to want to support pointers by dereferencing them and caching that as the
parameter but I've talked myself out of it. When checking a func's cache, should
the pointer point to the EXACT same object or just an object with equal values
(i.e. `reflect.DeepEqual(...)`)?
What if that struct has fields that are pointers? How could I verify circular
references are equal? What about nils? Typed vs Untyped Nils? These questions
are probably why pointers were never made hashable for use as map keys.

## Note to self:

```
gitsem <major|minor|patch>
git push
goreleaser release --rm-dist
```
