package parser

import (
	"github.com/amupitan/hero/ast/core"
	lx "github.com/amupitan/hero/lexer"
	"github.com/amupitan/hero/types"
)

type parser func(p *Parser) core.Expression

type Parser struct {
	// deprecated
	*lx.Lexer
	current lx.Token
	curr    int
	tokens  []lx.Token
	err     error
}

type CustomType string

var precedence = map[lx.TokenType]int{
	lx.Assign: 1, lx.Increment: 1, lx.Decrement: 1,
	lx.Or:       4,
	lx.And:      5,
	lx.LessThan: 9, lx.GreaterThan: 9, lx.LessThanOrEqual: 9, lx.GreaterThanOrEqual: 9, lx.Equal: 9, lx.NotEqual: 9,
	lx.Plus: 12, lx.Minus: 12,
	lx.Times: 15, lx.Div: 15, lx.Mod: 15,
}

var VALUES = []lx.TokenType{
	lx.Identifier,
	lx.Bool,
	lx.Int,
	lx.Float,
	lx.String,
	lx.RawString,
	lx.Rune,
	lx.Underscore,
}

var literals = VALUES[1:]

var builtins = map[string]types.Type{
	`bool`:    types.Bool,
	`float`:   types.Float,
	`func`:    types.Func,
	`generic`: types.Generic,
	`int`:     types.Int,
	`rune`:    types.Rune,
	`string`:  types.String,
}

// New returns a new parser
func New(input string) *Parser {
	p := &Parser{
		Lexer: lx.New(input),
	}

	p.tokens, p.err = p.Tokenize()
	return p
}

func (p *Parser) peek() *lx.Token {
	if p.curr >= len(p.tokens) {
		return nil
	}
	return &p.tokens[p.curr]
}

func (p *Parser) lookahead() *lx.Token {
	if p.curr+1 >= len(p.tokens) {
		return nil
	}
	return &p.tokens[p.curr+1]
}

func (p *Parser) next() *lx.Token {
	t := p.peek()
	if t != nil {
		p.curr++
	}
	return t
}

func (p *Parser) unstep() {
	if p.curr > 0 {
		p.curr--
	}
}

func (p *Parser) Parse() *core.Runtime {
	//TODO(DEV) parse imports
	return &core.Runtime{
		Body: p.parse_toplevel(),
	}
}

// delimited parses the content with a [start] and a [stop] token using the [separator]
// as a delimeter and [expr_parser] to parse the content. [end_sep] flags whether a separator is
// allowed at the end
// It uses an [p.parse_expression] if none is provided
func (p *Parser) delimited(start, stop, separator lx.TokenType, end_sep bool, expr_parser parser) []core.Expression {
	if !p.nextIs(start) {
		return nil
	}

	// consume start token
	p.expect(start)

	// if there is nothing between start and stop then exit
	if p.accept(stop) {
		p.next()
		return []core.Expression{}
	}

	// use expression parser if no parser is provided
	if expr_parser == nil {
		expr_parser = func(p *Parser) core.Expression { return p.parse_expression() }
	}

	params := make([]core.Expression, 0, 10) // TODO(CLEAN) we assume delimted content is usually ≤ 10

	// consume first expression before the delimeter
	params = append(params, expr_parser(p))

	// check for more content
	for {

		// consume and break when we see the stop token
		if p.accept(stop) {
			p.next()
			break
		}

		// consume separator
		p.expect(separator)

		if end_sep {
			// consume and break when we see the stop token
			if p.accept(stop) {
				p.next()
				break
			}
		}

		// consume expression
		params = append(params, expr_parser(p))
	}

	return params
}

// expect returns a toekn if the [expected] toekn type is the next token
// otherwise it panics
func (p *Parser) expect(expected lx.TokenType) *lx.Token {
	if expected != lx.NewLine {
		p.skipNewLines()
	}
	t := p.peek()

	if expected != lx.EndOfInput && t.Type == lx.EndOfInput {
		panic(p.reportEndOfInput(&expected))
	}

	if expected != t.Type {
		panic(p.reportUnexpected(&expected))
	}

	p.next()

	return t
}

// accept returns true if the [expected] token type matches the next token type
func (p *Parser) accept(expected lx.TokenType) bool {
	if expected != lx.NewLine {
		p.skipNewLines()
	}
	return p.nextIs(expected)
}

// nextIs returns true if the [expected] token type matches the next token type
func (p *Parser) nextIs(expected lx.TokenType) bool {
	if t := p.peek(); t != nil {

		return t.Type == expected
	}
	return false
}

// acceptsOneOf returns true if one of the expected tokens match
func (p *Parser) acceptsOneOf(expected ...lx.TokenType) bool {
	for i, _ := range expected {
		if p.accept(expected[i]) {
			return true
		}
	}
	return false
}

// expectsOneOf returns the first token that matches the current type or panics
// if none is found
func (p *Parser) expectsOneOf(expected ...lx.TokenType) *lx.Token {
	for i, _ := range expected {
		if p.accept(expected[i]) {
			return p.next()
		}
	}
	panic(p.reportUnexpectedMultiple(expected...))
}

// skipNewLines skips all new Line tokens till the next non-new Line token or the end
func (p *Parser) skipNewLines() {
	for p.peek().Type == lx.NewLine {
		p.next()
	}
}

func (c CustomType) IsType(value string) bool {
	return true
}

func (c CustomType) String() string {
	return string(c)
}
