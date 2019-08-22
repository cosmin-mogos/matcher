package matcher

import (
	"com.careem/matcher/parser"
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"reflect"
	"strconv"
)

type calcVisitor struct {
	*parser.BaseMatcherVisitor
	stack   valueStack
	failed  bool
	inArray bool
}

func (v *calcVisitor) VisitJson(ctx *parser.JsonContext) interface{} {
	return v.VisitValue(ctx.Value().GetRuleContext().(*parser.ValueContext))
}

func (v *calcVisitor) VisitObj(ctx *parser.ObjContext) interface{} {
	pairs := ctx.AllPair()
	matchers := make([]pairMatcher, len(pairs))
	for i, pair := range pairs {
		matchers[i] = v.VisitPair(pair.GetRuleContext().(*parser.PairContext)).(pairMatcher)
	}

	return objectMatcher{
		matchers: matchers,
	}
}

func (v *calcVisitor) VisitPair(ctx *parser.PairContext) interface{} {
	return pairMatcher{
		name:    stripQuotes(ctx.STRING()),
		matcher: v.VisitElement(ctx.Element().(*parser.ElementContext)).(Matcher),
	}
}

func (v *calcVisitor) VisitElement(ctx *parser.ElementContext) interface{} {
	switch {
	case ctx.Value() != nil:
		return v.VisitValue(ctx.Value().(*parser.ValueContext))
	case ctx.Template() != nil:
		return v.VisitTemplate(ctx.Template().(*parser.TemplateContext))
	default:
		panic("NOT IMPLEMENTED")
	}
}

func (v *calcVisitor) VisitArray(ctx *parser.ArrayContext) interface{} {
	elements := ctx.AllElement()
	matchers := make([]Matcher, len(elements))
	for i, element := range elements {
		matchers[i] = v.VisitElement(element.GetRuleContext().(*parser.ElementContext)).(Matcher)
	}

	return arrayMatcher{
		matchers: matchers,
	}
}

func (v *calcVisitor) VisitValue(ctx *parser.ValueContext) interface{} {
	switch {
	case ctx.Obj() != nil:
		return v.VisitObj(ctx.Obj().(*parser.ObjContext))
	case ctx.Array() != nil:
		return v.VisitArray(ctx.Array().(*parser.ArrayContext))
	case ctx.NUMBER() != nil:
		return Wrap(func(v interface{}) bool {
			//FIXME handle all int types
			rv := v.(reflect.Value)

			expected, err := strconv.ParseInt(ctx.NUMBER().GetText(), 10, 64)
			if err != nil {
				panic(fmt.Errorf("failed to parse \"%s\": %v", ctx.NUMBER().GetText(), err))
			}

			return rv.Type().Kind() == reflect.Int && rv.Int() == expected
		})
	default:
		panic("NOT IMPLEMENTED")
	}
}

func (v *calcVisitor) VisitTemplate(ctx *parser.TemplateContext) interface{} {
	switch {
	case ctx.ANY() != nil:
		return Wrap(func(v interface{}) bool {
			rv := v.(reflect.Value)
			return rv.IsValid()
		})
	case ctx.ANY_OR_OMIT() != nil:
		return Wrap(func(v interface{}) bool {
			//rv := v.(reflect.Value)
			//return rv.IsValid() || rv.IsNil()
			return true
		})
	default:
		panic("NOT IMPLEMENTED")
	}
}

func (v *calcVisitor) Visit(tree antlr.ParseTree) interface{}            { return nil }
func (v *calcVisitor) VisitChildren(node antlr.RuleNode) interface{}     { return nil }
func (v *calcVisitor) VisitTerminal(node antlr.TerminalNode) interface{} { return nil }
func (v *calcVisitor) VisitErrorNode(node antlr.ErrorNode) interface{}   { return nil }
