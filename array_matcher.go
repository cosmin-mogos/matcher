package matcher

import (
	"reflect"
)

type arrayMatcher struct {
	matchers []Matcher
}

func (m arrayMatcher) Matches(v interface{}) bool {
	rValue, ok := v.(reflect.Value)
	if !ok {
		rValue = reflect.ValueOf(v)
	}

	if rValue.Kind() != reflect.Slice && rValue.Kind() != reflect.Array {
		return false
	}

	if rValue.Len() != len(m.matchers) {
		return false
	}

	matches := true
	for i := 0; i < rValue.Len(); i++ {
		rElem := rValue.Index(i)

		matches = matches && m.matchers[i].Matches(rElem)
	}

	return matches
}
