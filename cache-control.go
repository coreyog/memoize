package memoize

import (
	"reflect"
)

type CacheControl interface {
	Clear()
	Remove(args ...interface{})
}

type cacheControl struct {
	cacheRoot map[interface{}]interface{}
}

func (cc *cacheControl) Clear() {
	cc.cacheRoot = map[interface{}]interface{}{}
}

func (cc *cacheControl) Remove(args ...interface{}) {
	var m, next interface{}
	var prev map[interface{}]interface{}
	var ok bool
	m = cc.cacheRoot

	for _, arg := range args {
		value := normalize(arg)

		switch mm := m.(type) {
		case map[interface{}]interface{}:
			next, ok = mm[value]
			if !ok {
				// haven't cached this set of args yet
				return
			}
		case []reflect.Value:
			// more args were provided than the memoized function has
			return
		}

		prev = m.(map[interface{}]interface{})
		m = next
	}

	normalLastArg := normalize(args[len(args)-1])
	leaf := prev[normalLastArg]

	switch l := leaf.(type) {
	case map[interface{}]interface{}:
		for key := range l {
			delete(l, key)
		}
	case []reflect.Value:
		prev[normalLastArg] = nil
	}
}

func (cc *cacheControl) fillAndCheck(args []reflect.Value) (results []reflect.Value) {
	var m, next interface{}
	var ok bool
	m = cc.cacheRoot

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

func (cc *cacheControl) fillAndSet(args []reflect.Value, results []reflect.Value) {
	var m interface{}
	var prev map[interface{}]interface{}

	m = cc.cacheRoot

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
