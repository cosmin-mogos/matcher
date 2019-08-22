package matcher

import (
	"reflect"
)

type valueStack []reflect.Value

func (s valueStack) Push(v reflect.Value) valueStack {
	return append(s, v)
}

func (s valueStack) Pop() valueStack {
	l := len(s)

	if l == 0 {
		return s
	}

	return s[:l-1]
}

func (s valueStack) Top() *reflect.Value {
	l := len(s)
	return &s[l-1]
}
