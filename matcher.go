package matcher

import (
	"com.careem/matcher/parser"
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"reflect"
	"strconv"
)

type calcListener struct {
	*parser.BaseMatcherListener
	v       valueStack
	failed  bool
	inArray bool
}

func (cl *calcListener) EnterValue(c *parser.ValueContext) {
	//top := cl.stack.Top()
	//fmt.Println(top)
	//fmt.Println(top.Type().Kind() == reflect.Int)
	switch {
	case c.NUMBER() != nil:
		//FIXME handle all number types
		expected, err := strconv.ParseInt(c.NUMBER().GetText(), 10, 64)
		if err == nil {
			cl.failed = cl.failed || cl.v.Top().Type().Kind() != reflect.Int || cl.v.Top().Int() != expected
		}

		expectedFloat, err := strconv.ParseFloat(c.NUMBER().GetText(), 64)
		if err == nil {
			cl.failed = cl.failed || cl.v.Top().Type().Kind() != reflect.Int || cl.v.Top().Float() != expectedFloat
		}

		if err != nil {
			panic("Could not handle NUMBER")
		}
		//fmt.Println("exp", expected, "|", c.NUMBER().GetText())
		//fmt.Println("real", cl.stack.Top().Int())

		//fmt.Println(cl.failed)
	case c.STRING() != nil:
		expected := stripQuotes(c.STRING())
		cl.failed = cl.failed || cl.v.Top().Type().Kind() != reflect.String || cl.v.Top().String() != expected
	}
}

func (cl *calcListener) EnterTemplate(c *parser.TemplateContext) {
	fmt.Println("Enter Template", c.GetText())
	switch {
	case c.ANY() != nil:
		// FIXME handle nil value properly
		cl.failed = cl.failed || (!cl.v.Top().IsValid()) // || cl.stack.Top().IsNil())
	case c.ANY_OR_OMIT() != nil:
	}
}

func (cl *calcListener) EnterArray(c *parser.ArrayContext) {
	fmt.Println("Enter Array", c.GetText())
	cl.inArray = true
}

func (cl *calcListener) ExitArray(c *parser.ArrayContext) {
	cl.inArray = false
}

func (cl *calcListener) EnterPair(c *parser.PairContext) {
	//fmt.Println("enter pair", c.STRING())

	fieldName := c.STRING().GetText()[1 : len(c.STRING().GetText())-1]
	fieldValue := cl.v.Top().FieldByName(fieldName)
	cl.v = cl.v.Push(fieldValue)

	//fmt.Println("lala", fieldValue.Int())
}

func (cl *calcListener) ExitPair(c *parser.PairContext) {
	//fmt.Println("exit pair", c.STRING())
	cl.v.Pop()
}

func match(template string, val interface{}) bool {
	// Setup the input
	is := antlr.NewInputStream(template)

	// Create the Lexer
	lexer := parser.NewMatcherLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := parser.NewMatcherParser(stream)

	// Finally parse the expression
	visitor := &calcVisitor{
		stack: []reflect.Value{reflect.ValueOf(val)},
	}
	matcher := p.Json().Accept(visitor).(Matcher)
	return matcher.Matches(val)
	//listener := &calcListener{
	//	stack: []reflect.Value{reflect.ValueOf(val)},
	//}
	//antlr.ParseTreeWalkerDefault.Walk(listener, p.Json())
	//return !listener.failed
}

func stripQuotes(strNode antlr.TerminalNode) string {
	text := strNode.GetText()
	return text[1 : len(text)-1]
}

type Matcher interface {
	Matches(v interface{}) bool
}
