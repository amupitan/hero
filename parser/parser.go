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
	lx.Assign:   1,
	lx.Or:       2,
	lx.And:      3,
	lx.LessThan: 7, lx.GreaterThan: 7, lx.LessThanOrEqual: 7, lx.GreaterThanOrEqual: 7, lx.Equal: 7, lx.NotEqual: 7,
	lx.Plus: 10, lx.Minus: 10,
	lx.Times: 20, lx.Div: 20, lx.Mod: 20,
}

var literals = []lx.TokenType{
	lx.Int,
	lx.Float,
	lx.String,
	lx.RawString,
	lx.Rune,
	lx.Underscore,
}

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
func (p *Parser) delimited(start, stop, separator lx.TokenType, end_sep bool, expr_parser parser) []core.Expression {
	if !p.accept(start) {
		return nil
	}

	// consume start token
	p.expect(start)

	// if there is nothing between start and stop then exit
	if p.accept(stop) {
		p.next()
		return []core.Expression{}
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
	t := p.peek()

	if expected != lx.EndOfInput && t.Type == lx.EndOfInput {
		return false
	}

	if expected != t.Type {
		return false
	}

	return true
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
