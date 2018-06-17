package parser

import (
	"errors"

	"github.com/amupitan/hero/ast"
	"github.com/amupitan/hero/ast/core"
	lx "github.com/amupitan/hero/lexer"
)

// parse_toplevel parses out the body of the program
func (p *Parser) parse_toplevel() core.Statement {
	var statements []core.Statement
	for t := p.peek(); t != nil && t.Type != lx.Unknown; {
		statements = append(statements, p.parse_statement())
	}
	return &ast.Program{Statements: statements}
}

// parse_statement parses a statement. It can parse any statement
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
	return p.parse_expression()
}

// parse_expression can parse any expression
func (p *Parser) parse_expression() core.Expression {
	return p.parse_binary(p.parse_atom(), nil)
}

// attempt_parse_call attempts to parse a call or returns nil if a call can't be parsed
func (p *Parser) attempt_parse_call() *ast.Call {
	object := ``
	identifier := p.expect(lx.Identifier)
	if p.accept(lx.Dot) {
		// consume dot
		p.next()

		object = identifier.Value
		identifier = p.expect(lx.Identifier)
	}
	params := p.delimited(lx.LeftParenthesis, lx.RightParenthesis, lx.Comma, false, func(p *Parser) core.Expression { return p.parse_expression() }) //TODO(CLEAN) parser arg
	if params == nil {
		// if parse was unsuccessful, retract and return
		p.unstep()
		return nil
	}

	// TODO: convert expression to call.params?
	return &ast.Call{
		Name:   identifier.Value,
		Args:   params,
		Object: object,
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
			// consume identifier as name
			name = p.next().Value

			// consume operator
			p.next()

			// get value
			value = p.parse_atom()
		} else {
			// return nil if a definition cannot be parsed
			return nil
		}
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

// parse_atom parses out an atom - which is a literal value or identifier
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

	// check if it is a call
	if p.accept(lx.Identifier) {
		if e := p.attempt_parse_call(); e != nil {
			return e
		}
	}

	t := p.expectsOneOf(lx.Identifier,
		lx.Int, lx.Float,
		lx.String,
		lx.RawString,
		lx.Rune,
		lx.Underscore)

	// TODO: allow functions
	return &ast.Atom{
		Type:  t.Type,
		Value: t.Value,
	}
}

// parse_binary parses a binary expression
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

	e := p.parse_binary(b, my_op)

	if op.Type == lx.Assign {
		return p.parse_assignment(e)
	}

	return e
}

// parse_assignment parses an assigment expression
func (p *Parser) parse_assignment(e core.Expression) core.Expression {
	// TODO(DEV) check that b.Left is an identifier
	if b, ok := e.(*ast.Binary); ok {
		return &ast.Assignment{
			Identifier: b.Left.String(),
			Value:      b.Right,
		}
	}
	// TODO(DEV) find a better way to take care of invalid states
	panic(errors.New(`Cannot assign variable to an assignment`))
}
