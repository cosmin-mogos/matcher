package matcher

import "reflect"

type objectMatcher struct {
	matchers []pairMatcher
}

func (m objectMatcher) Matches(v interface{}) bool {
	parentValue, ok := v.(reflect.Value)
	if !ok {
		parentValue = reflect.ValueOf(v)
	}

	matches := true
	for _, matcher := range m.matchers {
		matches = matches && matcher.Matches(parentValue)
	}

	return matches
}
