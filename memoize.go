package memoize

import (
	"fmt"
	"reflect"
)

func Memo[T any](fn T) (m T, err error) {
	fnt := reflect.TypeOf(fn)

	if fnt.Kind() != reflect.Func {
		return fn, ErrNotAFunc
	}

	if fnt.NumIn() == 0 {
		return fn, ErrMissingArgs
	}

	if fnt.NumOut() == 0 {
		return fn, ErrMissingReturns
	}

	fnv := reflect.ValueOf(fn)

	cacheRoot := map[interface{}]interface{}{}

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

		cResults := fillAndCheck(cacheRoot, args)
		if cResults == nil {
			if fnt.IsVariadic() {
				results = fnv.CallSlice(args)
			} else {
				results = fnv.Call(args)
			}

			fillAndSet(cacheRoot, args, results)
		} else {
			results = cResults
		}

		return results
	})

	return ret.Interface().(T), nil
}

func fillAndCheck(cacheRoot map[interface{}]interface{}, args []reflect.Value) (results []reflect.Value) {
	var m, next interface{}
	var ok bool
	m = cacheRoot

	for _, arg := range args {
		value := normalize(arg.Interface())

		next, ok = m.(map[interface{}]interface{})[value]
		if !ok {
			next = map[interface{}]interface{}{}
			m.(map[interface{}]interface{})[value] = next
		}

		m = next
	}

	results, ok = m.([]reflect.Value)

	if !ok || results == nil {
		return nil
	}

	return results
}

func fillAndSet(cacheRoot map[interface{}]interface{}, args []reflect.Value, results []reflect.Value) {
	var m interface{}
	var prev map[interface{}]interface{}

	m = cacheRoot

	for _, arg := range args {
		prev = m.(map[interface{}]interface{})
		value := normalize(arg.Interface())

		m = prev[value]
	}

	value := normalize(args[len(args)-1].Interface())

	prev[value] = results
}

func normalize(arg interface{}) (root interface{}) {
	argT := reflect.TypeOf(arg)
	if argT.Kind() == reflect.Slice {
		return normalizeSlice(arg)
	}

	return arg
}

func normalizeSlice(arg interface{}) (norm interface{}) {
	argV := reflect.ValueOf(arg)
	length := argV.Len()

	arr := reflect.New(reflect.ArrayOf(length, argV.Type().Elem())).Elem()

	for i := 0; i < length; i++ {
		arr.Index(i).Set(reflect.ValueOf(normalize(argV.Index(i).Interface())))
	}

	reflect.Copy(arr, argV)

	return arr.Interface()
}
