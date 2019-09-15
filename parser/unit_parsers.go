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

	p.skipNewLines()
	// TODO(CLEAN) remove check for unknown, it shouldn't get here
	for t := p.peek(); t != nil && t.Type != lx.Unknown && t.Type != lx.EndOfInput; t = p.peek() {
		statements = append(statements, p.parse_statement())
		p.skipNewLines()
	}
	return &ast.Program{Body: &ast.Block{
		Statements: statements,
	}}
}

// parse_statement parses a statement. It can parse any statement
func (p *Parser) parse_statement() core.Statement {
	t := p.peek()
	switch t.Type {
	case lx.LoopName:
		fallthrough
	case lx.For:
		return p.parse_loop()
	case lx.Func:
		return p.parse_func(false)
	case lx.If:
		return p.parse_if()
	case lx.LeftBrace:
		return p.parse_block()
	case lx.Return:
		return p.parse_return()
		//TODO

	}

	// attempt to parse definition
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
// it parses a lambda function if possible
func (p *Parser) attempt_parse_call() core.Expression {

	if p.accept(lx.Identifier) {
		// This can't simply be changed to return p.attempt_parse_named_call()
		// because p.attempt_parse_named_call() can never return nil since
		// it has a type. See: https://golang.org/doc/faq#nil_error
		if c := p.attempt_parse_named_call(); c != nil {
			return c
		}
		return nil
	}
	return p.attempt_parse_lambda_call()
}

// attempt_parse_lambda_call attempts to parse a call
// from a lambda expression, returns the lambda expression
// if it is not call or panics if neither is possible
func (p *Parser) attempt_parse_lambda_call() core.Expression {
	f := p.parse_func(true)
	if t := p.peek(); t.Type == lx.LeftParenthesis {
		return &ast.Call{
			Args: p.delimited(lx.LeftParenthesis, lx.RightParenthesis, lx.Comma, false, nil),
			Func: f,
		}
	}

	return f
}

// attempt_parse_named_call attempts to parse a call
// from an identifier
func (p *Parser) attempt_parse_named_call() *ast.Call {
	object := ``
	identifier := p.expect(lx.Identifier)
	if p.nextIs(lx.Dot) {
		// consume dot
		p.next()

		object = identifier.Value
		identifier = p.expect(lx.Identifier)
	}
	params := p.delimited(lx.LeftParenthesis, lx.RightParenthesis, lx.Comma, false, nil)
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
	var name, Type *lx.Token
	var value core.Expression
	if p.accept(lx.Var) {
		// consume var keyword
		p.next()

		// consume identifier name
		name = p.expect(lx.Identifier)

		// check if type is present
		if p.accept(lx.Identifier) {
			Type = p.next()

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

	} else if p.accept(lx.Identifier) {
		if lookahead := p.lookahead(); lookahead != nil && lookahead.Type == lx.Declare {
			// consume identifier as name
			name = p.next()

			// consume operator
			p.next()

			// get value
			value = p.parse_binary(p.parse_atom(), nil)
		} else {
			// return nil if a definition cannot be parsed
			return nil
		}
	} else {
		// return nil if a definition cannot be parsed
		return nil
	}

	if Type == nil {
		Type = &lx.Token{}
	}

	return &ast.Definition{
		Name:      *name,
		Value:     value,
		LexerType: *Type,
	}
}

// parse_atom parses out an atom - which is a literal value or identifier
func (p *Parser) parse_atom() core.Expression {
	isSigned := false
	// check for negation
	isNegated := p.accept(lx.Not)
	if isNegated {
		// consume negation token
		p.next()
	} else if p.acceptsOneOf(lx.Plus, lx.Minus) {
		// check for specified sign (+ or -)
		isSigned = p.nextIs(lx.Minus)
		// consume + or - token
		p.next()
	}

	signAndOrNegate := func(exp core.Expression) {
		// attempt to negate if there was a negation
		if isNegated {
			negateExpr(exp)
			return
		}

		// attempt to sign if there was a minus sign
		if isSigned {
			signExpr(exp)
		}
	}

	// attempt to consume expression in a parenthesis
	if p.accept(lx.LeftParenthesis) {
		// skip left paren
		p.next()
		exp := p.parse_expression()

		// consume right paren
		p.expect(lx.RightParenthesis)

		signAndOrNegate(exp)
		return exp
	}

	// parse call if it is a named or lambda call
	if p.nextIs(lx.Identifier) || p.nextIs(lx.Func) {
		if e := p.attempt_parse_call(); e != nil {
			signAndOrNegate(e)
			return e
		}
	}

	t := p.expectsOneOf(VALUES...)

	if isNegated && !isBooleanAble(t.Type) {
		// TODO(REPORT) better message
		report(`cannot negate non-boolean type`)
		return nil
	}

	if isSigned && !isSignSpecifiable(t.Type) {
		// TODO(REPORT) better message
		report(`cannot specify sign of non-number type`)
		return nil
	}

	// TODO: allow functions
	return &ast.Atom{
		Token:   *t,
		Negated: isNegated,
		Signed:  isSigned,
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

	if p.acceptsOneOf(lx.Increment, lx.Decrement) {
		return p.parse_assignment(left)
	}

	// consume operator
	op := p.next()

	right := p.parse_binary(p.parse_atom(), &(op.Type))
	b := &ast.Binary{
		Left:     left,
		Operator: *op,
		Right:    right,
	}

	// check for invalid boolean expressions
	if op.Type == lx.And || op.Type == lx.Or {
		ensureBoolean(left, right)
	}

	e := p.parse_binary(b, my_op)

	if isBinaryAssgnmentToken(op.Type) {
		return p.parse_assignment(e)
	}

	return e
}

// parse_assignment parses an assigment expression
// TODO(DEV) split into one that takes in an atom and another a binary
func (p *Parser) parse_assignment(e core.Expression) core.Expression {
	switch a := e.(type) {
	case *ast.Binary:
		if !isIdentifier(a.Left) {
			break
		}
		var value core.Expression

		// if it is a pure assignment, then the right side is the value
		if a.Operator.Type == lx.Assign {
			value = a.Right
		} else {
			// it is an operation assignment
			value = &ast.Operation{Token: a.Operator, Value: a.Right}
		}
		return &ast.Assignment{
			Identifier: a.Left.String(),
			Value:      value,
		}
	case *ast.Atom:
		if a.Token.Type == lx.Identifier {
			t := p.expectsOneOf(lx.Increment, lx.Decrement)
			return &ast.Assignment{
				Identifier: a.Value,
				Value:      &ast.Operation{Token: *t},
			}
		}
	}
	// TODO(DEV) find a better way to take care of invalid states
	report(`Cannot assign value to non-identifier`)

	// report panics so this will never be hit
	return nil
}

// isIdentifier returns true if an expression is an identifier
// TODO(DEV) there should be a Type() in core.Expression interface instead
func isIdentifier(e core.Expression) bool {
	if id, ok := e.(*ast.Atom); ok {
		return id.Token.Type == lx.Identifier
	}
	return false
}

// isAssgnmentToken returns true if the token type
// is an assignment that can be used in a binary expression
func isBinaryAssgnmentToken(t lx.TokenType) bool {
	return t == lx.Assign || t == lx.PlusEq || t == lx.MinusEq ||
		t == lx.TimesEq || t == lx.DivEq || t == lx.ModEq
}

// parse_block parses a block surrounded by braces
func (p *Parser) parse_block() *ast.Block {
	// consume left brace
	p.expect(lx.LeftBrace)

	// return an empty slice if there are no statements
	if p.accept(lx.RightBrace) {
		p.next()
		return &ast.Block{}
	}

	// we assume blocks are usually <= 20 statements
	statements := make([]core.Statement, 0, 20)
	for !p.accept(lx.RightBrace) {
		statements = append(statements, p.parse_statement())
		// TODO(DEV) expect semi-colon or new line? Consider one-liners
	}

	// consume right brace
	p.next()

	return &ast.Block{
		Statements: statements,
	}
}

// parse_func parses a function
func (p *Parser) parse_func(lamdba bool) *ast.Function {

	var name *lx.Token
	// consume func
	p.expect(lx.Func)

	// consume function name if not lambda
	if !lamdba {
		name = p.expect(lx.Identifier)
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
			Name: *name,
			// LexerType: types.Func.String(), // TODO(DEV) remove String() caller
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

// parse_if parses an if statement
func (p *Parser) parse_if() *ast.If {
	p.expect(lx.If)

	hasLeftParen := false
	// attempt to consume expression in a parenthesis
	if p.accept(lx.LeftParenthesis) {
		// skip left paren
		p.next()
		hasLeftParen = true
	}
	// get condition
	cond := p.parse_expression()

	// consume right paren if a left was consumed
	if hasLeftParen {
		p.expect(lx.RightParenthesis)
	}

	if !isBooleanExpr(cond) {
		// TODO(DEV) include line and column information
		report(`Only boolean expressions are allowed in if statements`)
		return nil
	}
	body := p.parse_block()

	// TODO(DEV) determine whether new lines are allowed between if elses

	var else_ *ast.If
	if p.accept(lx.Else) {
		// consume else token
		p.next()

		// check if it's an else-if
		if p.nextIs(lx.If) {
			else_ = p.parse_if()
		} else {
			else_ = &ast.If{
				Body: p.parse_block(),
			}
		}
	}

	return &ast.If{
		Condition: cond,
		Body:      body,
		Else:      else_,
	}
}

// parse_loop parses a for statement
func (p *Parser) parse_loop() ast.Loop {
	// check if loop is named
	var name string
	if p.accept(lx.LoopName) {
		name = p.next().Value

		// there was be a new line after a loop name
		p.expect(lx.NewLine)
	}

	if rl := p.attempt_parse_range_loop(); rl != nil {
		rl.Name = name
		return rl
	}
	// 'for' token is already consumed

	var preLoop, condition, postIter core.Expression

	// if loop has no statements then parse the body
	if p.nextIs(lx.LeftBrace) {
		return &ast.ForLoop{Name: name, Body: p.parse_block()}
	}

	// if a stement exists before the first semicolon
	// then it is the preLoop statement
	if !p.nextIs(lx.SemiColon) {
		if d := p.attempt_parse_definition(); d != nil {
			preLoop = d
		} else {
			preLoop = p.parse_expression()
		}
	}

	// if there's only one statement in the loop, then
	// it's a condition-only loop and we're done
	if !p.nextIs(lx.SemiColon) {
		return &ast.ForLoop{
			Name:      name,
			Condition: preLoop,
			Body:      p.parse_block(),
		}
	}

	// consume semicolon
	p.next()

	// check if there's a condition statement before the
	// next semicolon
	if !p.nextIs(lx.SemiColon) {
		// parse condition
		condition = p.parse_expression()
	}

	// consume semicolon
	p.expect(lx.SemiColon)

	// check if there's more tokens before the body's brace
	if !p.nextIs(lx.LeftBrace) {
		postIter = p.parse_expression()
	}

	// parse_body:
	return &ast.ForLoop{
		Name:          name,
		PreLoop:       preLoop,
		Condition:     condition,
		PostIteration: postIter,
		Body:          p.parse_block(),
	}
}

func (p *Parser) attempt_parse_range_loop() *ast.RangeLoop {
	p.expect(lx.For)
	// get parser cursor before parse attempt
	initial := p.curr

	// flag for successfully parsing a range loop
	success := false

	// the value of the second identifier for the range loop
	second := ``

	// restore parser cursor if range loop
	// was not successfully parsed
	defer func() {
		if !success {
			p.curr = initial
		}
	}()

	// consume first identifier
	if !p.accept(lx.Identifier) { // TODO(DEV) or underscore
		return nil
	}
	first := p.next().Value

	// potentially consume in
	if p.accept(lx.In) {
		goto parse_iterable
	}

	// consume comma
	if !p.nextIs(lx.Comma) {
		return nil
	}
	p.next()

	// consume second identifier
	if !p.nextIs(lx.Identifier) { // TODO(DEV) or underscore
		return nil
	}
	second = p.next().Value

	if !p.nextIs(lx.In) {
		return nil
	}

parse_iterable:
	// consume in token
	p.next()

	iterable := p.expect(lx.Identifier).Value
	success = true
	return &ast.RangeLoop{
		First:    first,
		Second:   second,
		Iterable: iterable,
		Body:     p.parse_block(),
	}
}

// parse_return parses a return statement
func (p *Parser) parse_return() *ast.Return {
	// consume return token
	p.expect(lx.Return)

	values := p.delimited(lx.LeftParenthesis, lx.RightParenthesis, lx.Comma, true, nil)
	if values != nil {
		return &ast.Return{Values: values}
	}

	values = make([]core.Expression, 0, 5) // we assume most return statements will have ≤ 5 values
	for {
		// get next return value
		isLiteral := func() bool {
			t := p.peek().Type
			for _, literal := range literals {
				if t == literal {
					return true
				}
			}
			return false
		}
		// TODO(DEV) use nextIs(...)
		if p.nextIs(lx.Identifier) || p.nextIs(lx.Func) || isLiteral() {
			values = append(values, p.parse_expression())
			// TODO(DEV) use a universal check for end of input
		} else if !p.nextIs(lx.EndOfInput) {
			// TODO(REPORT) use expectsOneOf or something better
			report(`Expected an expression`)
			return nil
		}

		// break if no comma is found
		if !p.nextIs(lx.Comma) {
			break
		}

		// consume comma
		p.next()
	}
	return &ast.Return{
		Values: values,
	}
}

// isBooleanAble returns true if the token could
// possibly be a boolean value
func isBooleanAble(t lx.TokenType) bool {
	return t == lx.Bool || t == lx.Identifier
}

// isSignSpecifiable returns true if the token could
// be sign specified
func isSignSpecifiable(t lx.TokenType) bool {
	return t == lx.Int || t == lx.Float || t == lx.Identifier
}

// isBooleanBinaryExpr returns true if the binary
// expression's operator is a comparator
// i.e. ==, <, >, <=, !=
func isBooleanBinaryExpr(op lx.TokenType) bool {
	return op == lx.Equal || op == lx.LessThan || op == lx.GreaterThan || op == lx.LessThanOrEqual ||
		op == lx.NotEqual || op == lx.GreaterThanOrEqual || op == lx.And || op == lx.Or
}

// isArithmeticBinaryExpr returns true if the binary
// expression's operator is arithmetic
// i.e. +, -, /, %
func isArithmeticBinaryExpr(op lx.TokenType) bool {
	return op == lx.Plus || op == lx.Minus || op == lx.Times || op == lx.Div || op == lx.Mod
}

// isBooleanExpr returns true if the expression is a
// valid boolean expression
// e.g. identifier, boolean binary expression & calls
// only boolean expressions can be negated
func isBooleanExpr(e core.Expression) bool {
	switch exp := e.(type) {
	case *ast.Atom:
		return exp.Token.Type == lx.Bool || exp.Token.Type == lx.Identifier
	case *ast.Binary:
		return isBooleanBinaryExpr(exp.Operator.Type)
	case *ast.Call:
		return true
	}

	return false
}

// ensureBoolean fails if one of the expressions
// is not a boolean expression
func ensureBoolean(exps ...core.Expression) {
	for _, e := range exps {
		if !isBooleanExpr(e) {
			report(e.String() + ` is used in a boolean context but is not a boolean expression`)
		}
	}
}

// negateExpr negates a booleanable expression
// TODO(DEV) this should be moved a booleanable interface
// that has a Negate() method
func negateExpr(e core.Expression) {
	switch exp := e.(type) {
	case *ast.Atom:
		exp.Negated = true
	case *ast.Binary:
		if isBooleanBinaryExpr(exp.Operator.Type) {
			exp.Negated = true
			return
		}
		report(`cannot negate non-boolean expression`)
	case *ast.Call:
		exp.Negated = true
	default:
		// TODO(REPORT) better message
		report(`cannot negate non-boolean expression`)
	}
}

// signExpr signs a negative number or call expression
// TODO(REPORT) better message in reports
func signExpr(e core.Expression) {
	switch exp := e.(type) {
	case *ast.Atom:
		exp.Signed = true
	case *ast.Call:
		exp.Signed = true
	case *ast.Binary:
		if isArithmeticBinaryExpr(exp.Operator.Type) {
			exp.Signed = true
			return
		}
		report(`cannot specify sign of non-number expression`)
	default:
		report(`cannot specify sign of non-number type`)
	}
}
