package matcher

import "reflect"

type pairMatcher struct {
	name    string
	matcher Matcher
}

func (m pairMatcher) Matches(v interface{}) bool {
	parentValue, ok := v.(reflect.Value)
	if !ok {
		parentValue = reflect.ValueOf(v)
	}

	fieldValue := parentValue.FieldByName(m.name)
	return m.matcher.Matches(fieldValue)
}
