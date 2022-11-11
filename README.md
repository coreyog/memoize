# Memoize

A generic memoization library.

Needs more testing before going to prod.

See tests for examples.

Not every function can be memoized. The types of allowed parameters are limited
by this library. I believe the return can be anything (channels, funcs, maps).

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

Because I want to return a func with the exact same parameter list and return values
as the provided input, any errors that arise will panic.
