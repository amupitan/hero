package parser

import (
	"github.com/amupitan/hero/ast"
	"github.com/amupitan/hero/ast/core"
	lx "github.com/amupitan/hero/lexer"
	"github.com/amupitan/hero/types"
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
	case lx.LeftBrace:
		//TODO

	}

	// attempt to parse short definition
	if d := p.attempt_parse_definition(); d != nil {
		return d
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
				value = p.parse_binary(p.parse_atom(), nil)
			}
		} else {
			// if type isn't present, then there must be a value
			// consume assignment token
			p.expect(lx.Assign)

			value = p.parse_binary(p.parse_atom(), nil)
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
			value = p.parse_binary(p.parse_atom(), nil)
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
	report(`Cannot assign variable to an assignment`)

	// report panics so this will never be hit
	return nil
}

// parse_block parses a block surrounded by braces
func (p *Parser) parse_block() []core.Statement {
	// consume left brace
	p.expect(lx.LeftBrace)

	// return an empty slice if there are no statements
	if p.accept(lx.RightBrace) {
		p.next()
		return []core.Statement{}
	}

	// we assume blocks are usually <= 20 statements
	statements := make([]core.Statement, 0, 20)
	for !p.accept(lx.RightBrace) {
		statements = append(statements, p.parse_statement())
		// TODO(DEV) expect semi-colon or new line? Consider one-liners
	}

	// consume right brace
	p.next()

	return statements
}

// parse_func parses a function
func (p *Parser) parse_func(lamdba bool) *ast.Function {

	var name string
	// consume func
	p.expect(lx.Func)

	// consume function name if not lambda
	if !lamdba {
		name = p.expect(lx.Identifier).Value
	}

	// get function parameters
	params := p.parse_func_params()

	// we assume most functions have returns ≤ 5
	returns := make([]types.Type, 0, 5)

	getType := func(identifier string) types.Type {
		// check if type is a builtin else
		// create custom type
		if _type, ok := builtins[identifier]; ok {
			return _type
		}
		return CustomType(identifier)
	}

	// get return types
	//
	// has one return type
	if p.accept(lx.Identifier) {
		_type := getType(p.next().Value)
		returns = append(returns, _type)
	} else if p.accept(lx.LeftParenthesis) {
		rets := p.delimited(lx.LeftParenthesis, lx.RightParenthesis, lx.Comma, false, func(p *Parser) core.Expression {
			t := p.expect(lx.Identifier)
			return &ast.Value{Value: t.Value}
		})

		// add parsed return types
		for i := range rets {
			name := rets[i].(*ast.Value).Value
			returns = append(returns, getType(name))
		}
	}

	// parse function body
	body := p.parse_block()

	return &ast.Function{
		Definition: ast.Definition{
			Name: name,
			Type: types.Func.String(), // TODO(DEV) remove String() caller
		},
		Parameters:  params,
		Body:        body,
		ReturnTypes: returns,
		Lambda:      lamdba,
	}
}

// parse_func_params parses the parameters from a function
func (p *Parser) parse_func_params() []*ast.Param {

	// consume left paren
	p.expect(lx.LeftParenthesis)

	// if there is nothing between the parenthesis,
	// return an empty slice
	if p.accept(lx.RightParenthesis) {
		p.next()
		return []*ast.Param{}
	}

	// we assume most functions have params ≤ 10
	params := make([]*ast.Param, 0, 10)

	// buffer to store identifier names till their
	// type has been identified
	buff := make([]string, 0, 5)

	for {
		// get next parameter name
		identifier := p.expect(lx.Identifier).Value

		// add parameter name to buffer
		buff = append(buff, identifier)

		// if type is founf
		if p.accept(lx.Identifier) {
			var (
				_type types.Type
				ok    bool
			)
			// get type name
			typeName := p.next().Value

			// check if type is a builtin else
			// create custom type
			if _type, ok = builtins[typeName]; !ok {
				_type = CustomType(typeName)
			}

			// create a param from everything in the buffer
			// and assign the type that was found to each of those
			// params created
			for i := range buff {
				param := &ast.Param{Name: buff[i], Type: _type}
				params = append(params, param)
			}

			// empty the buffer
			buff = buff[:0]

			// check for comma separator or end parenthesis
			next := p.expectsOneOf(lx.Comma, lx.RightParenthesis)
			if next.Type == lx.Comma {
				continue
			}

			if next.Type == lx.RightParenthesis {
				break
			}
		}

		// check for comma separator
		if p.accept(lx.Comma) {
			p.next()
			continue
		}

		// panic for type not found
		//
		// if it got here then there is something in the
		// buffer that hasn't been added to the params because
		// no typename was found
		p.expect(lx.Identifier)
	}

	return params

}
