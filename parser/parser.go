package parser

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/amupitan/hero/ast"
	"github.com/amupitan/hero/ast/core"
	lx "github.com/amupitan/hero/lexer"
)

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

type parser func(p *Parser) core.Expression

type Parser struct {
	// deprecated
	*lx.Lexer
	current lx.Token
	curr    int
	tokens  []lx.Token
	err     error
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

// parse_toplevel parses out the body of the program
func (p *Parser) parse_toplevel() core.Statement {
	var statements []core.Statement
	for t := p.peek(); t != nil && t.Type != lx.Unknown; {
		statements = append(statements, p.parse_statement())
	}
	return &ast.Program{Statements: statements}
}

func (p *Parser) parse_statement() core.Statement {
	t := p.peek()
	switch t.Type {
	case lx.Var:
		return p.attempt_parse_definition()
	case lx.Return:
	case lx.If:
	case lx.For:
		//TODO

	}
	return &ast.Function{}
}

func (p *Parser) parse_expression() core.Expression {
	return p.parse_binary(p.parse_atom(), nil)
}

// attempt_parse_call attempts to parse a call or returns nil if a call can't be parsed
func (p *Parser) attempt_parse_call() *ast.Call {
	identifier := p.expect(lx.Identifier)
	params := p.delimeted(lx.LeftParenthesis, lx.RightParenthesis, lx.Comma, func(p *Parser) core.Expression { return p.parse_expression() }) //TODO(CLEAN) parser arg
	if params == nil {
		// if parse was unsuccessful, retract and return
		p.unstep()
		return nil
	}

	// TODO: convert expression to call.params?
	return &ast.Call{
		Name: identifier.Value,
		Args: params,
	}
}

// attempt_parse_definition attempts to parse a definition or returns nil if it can't be parsed
func (p *Parser) attempt_parse_definition() *ast.Definition {
	var name, Type string
	var value core.Expression
	if p.accept(lx.Var) {
		// consume var keyword
		p.next()

		// consume identifier name
		name = p.expect(lx.Identifier).Value

		// check if type is present
		if p.accept(lx.Identifier) {
			Type = p.next().Value

			// consume value if assign token is present
			if p.accept(lx.Assign) {
				p.next()
				// get value
				value = p.parse_atom()
			}
		} else {
			// if type isn't present, then there must be a value
			// cosume assigment token
			p.expect(lx.Assign)

			value = p.parse_atom()
		}

	} else if p.accept(lx.Func) {
		// TODO: parse func
	} else if p.accept(lx.Identifier) {
		if lookahead := p.lookahead(); lookahead != nil && lookahead.Type == lx.Declare {
			// cosume identifier as name
			name = p.next().Value

			// get value
			value = p.parse_atom()
		}

		// return nil if a definition cannot be parsed
		return nil
	} else {
		// panic if token is invalid
		//
		//TODO(IMPROV) use a function that will
		// straight up panic instead of an unnecessary loop
		p.expectsOneOf(lx.Var, lx.Identifier, lx.Func)
	}

	return &ast.Definition{
		Name:  name,
		Value: value,
		Type:  Type,
	}
}

func (p *Parser) parse_atom() core.Expression {
	// attempt to consume expression in a parenthesis
	if p.accept(lx.LeftParenthesis) {
		// skip left paren
		p.next()
		exp := p.parse_expression()

		// consume right paren
		p.expect(lx.RightParenthesis)
		return exp
	}

	t := p.expectsOneOf(lx.Identifier,
		lx.Int, lx.Float,
		lx.String,
		lx.RawString,
		lx.Rune,
		lx.Underscore)

	// check if it is a call
	if t.Type == lx.Identifier {
		if e := p.attempt_parse_call(); e != nil {
			return e
		}
	}

	// TODO: allow functions
	return &ast.Atom{
		Type:  t.Type,
		Value: t.Value,
	}
}

func (p *Parser) parse_binary(left core.Expression, my_op *lx.TokenType) core.Expression {
	var (
		prec int
		ok   bool
	)

	if my_op != nil {
		if prec, ok = precedence[*my_op]; !ok {
			return left
		}
	}

	// return left if next token isn't an operation
	// or the next precedence is lower
	if op := p.peek(); op != nil {
		if next_prec, ok := precedence[op.Type]; !ok || prec >= next_prec {
			return left
		}
	}

	// consume operator
	op := p.next()

	right := p.parse_binary(p.parse_atom(), &(op.Type))
	b := &ast.Binary{
		Left:     left,
		Operator: *op,
		Right:    right,
	}

	return p.parse_binary(b, my_op)
}

func (p *Parser) delimeted(start, stop, separator lx.TokenType, expr_parser parser) []core.Expression {
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

	params := make([]core.Expression, 0, 10) // TODO(CLEAN) we assume delimted content is usually < 10

	// consume first expression before the delimeter
	params = append(params, expr_parser(p))

	// check for more content
	for {
		// consume separator
		p.expect(separator)

		// consume expression
		params = append(params, expr_parser(p))

		// consume and break when we see the stop token
		if p.accept(stop) {
			p.next()
			break
		}
	}

	return params
}

func (p *Parser) parse_assignment() core.Expression {
	identifier := p.expect(lx.Identifier)

	// expect and ignore assignment token
	p.expect(lx.Assign)

	value := p.parse_expression()
	return &ast.Assignment{
		Identifier: *identifier,
		Value:      value,
	}

}

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

func (p *Parser) reportEndOfInput(expected *lx.TokenType) error {
	// TODO(DEV) add file name
	return fmt.Errorf("%d:%d: Expected `%s` but reached end of file.", p.Lexer.Line, p.Lexer.Column, *expected)
}

func (p *Parser) reportUnexpected(expected *lx.TokenType) error {
	// TODO(DEV) add file name
	return fmt.Errorf("%d:%d: Expected `%v` but found `%s`.", p.Lexer.Line, p.Lexer.Column, *expected, p.peek().Value)
}

func (p *Parser) reportUnexpectedMultiple(expected ...lx.TokenType) error {
	// TODO(DEV) add file name
	sb := bytes.Buffer{}

	// Add line and column info
	sb.WriteRune(rune(p.Lexer.Line)) // TODO(DEV) this or rune-int calculation
	sb.WriteRune(':')
	sb.WriteRune(rune(p.Lexer.Column))
	sb.WriteString(": Expected either ")
	for _, ex := range expected {
		sb.WriteString(string(ex))
		sb.WriteString(", ")
	}

	// remove last comma and space
	sb.Truncate(sb.Len() - 2)

	sb.WriteString("but received ")
	if curr := p.peek(); curr != nil {
		sb.WriteString(curr.Value)
	} else {
		// TODO: should never get here
		sb.WriteString("Unknown")
	}

	return errors.New(sb.String())
}
