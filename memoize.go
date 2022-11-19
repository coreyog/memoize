package memoize

import (
	"fmt"
	"reflect"
)

func Memo[T any](fn T) (m T, cache CacheControl, err error) {
	fnt := reflect.TypeOf(fn)

	if fnt.Kind() != reflect.Func {
		return fn, nil, ErrNotAFunc
	}

	if fnt.NumIn() == 0 {
		return fn, nil, ErrMissingArgs
	}

	if fnt.NumOut() == 0 {
		return fn, nil, ErrMissingReturns
	}

	fnv := reflect.ValueOf(fn)

	cc := &cacheControl{
		cacheRoot: map[interface{}]interface{}{},
	}

	ret := reflect.MakeFunc(fnt, func(args []reflect.Value) (results []reflect.Value) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}

			err, ok := r.(error)
			if !ok {
				// should always be an error, but just in case...
				panic(r)
			}

			panic(fmt.Errorf("panic in memo stub: %w", err))
		}()

		cResults := cc.fillAndCheck(args)
		if cResults == nil {
			if fnt.IsVariadic() {
				results = fnv.CallSlice(args)
			} else {
				results = fnv.Call(args)
			}

			cc.fillAndSet(args, results)
		} else {
			results = cResults
		}

		return results
	})

	return ret.Interface().(T), cc, nil
}
