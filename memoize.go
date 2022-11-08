package memoize

import "reflect"

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
		cResults := fillAndCheck(cacheRoot, args)
		if cResults == nil {
			results = fnv.Call(args)
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
		next, ok = m.(map[interface{}]interface{})[arg.Interface()]
		if !ok {
			next = map[interface{}]interface{}{}
			m.(map[interface{}]interface{})[arg.Interface()] = next
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
		m = prev[arg.Interface()]
	}

	prev[args[len(args)-1].Interface()] = results
}
