package matcher

type MatcherFunc func(v interface{}) bool

type wrapperMatcher struct {
	f MatcherFunc
}

func (m wrapperMatcher) Matches(v interface{}) bool {
	return m.f(v)
}

func Wrap(f MatcherFunc) Matcher {
	return wrapperMatcher{f: f}
}
