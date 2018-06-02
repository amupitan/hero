package lexer

import (
	"bytes"
)

// Lexer performs lexical analysis on an input
type Lexer struct {
	input                  []byte
	position, line, column int
}

func New(input string) *Lexer {
	return &Lexer{
		input: []byte(input),
	}
}

func (l *Lexer) NextToken() Token {
	if l.position >= len(l.input) {
		return EndOfInputToken
	}

	curr := l.getCurr()

	if isParenthesis(curr) {
		return l.consumeParenthesis()
	}

	return UnknownToken
}

func isLetter(b byte) bool              { return false }
func isDigit(b byte) bool               { return false }
func isValidIdentifierChar(b byte) bool { return false }
func isOperator(b byte) bool            { return false }
func isParenthesis(b byte) bool         { return false }
func isArithmeticOperator(b byte) bool  { return false }
func isComparisonOperator(b byte) bool  { return false }

// getCurr returns the byte at the current position
func (l *Lexer) getCurr() byte {
	return l.input[l.position]
}

// moves the lexer a position forward
func (l *Lexer) move() {
	l.position++
	l.column++
}

// consumeParenthesis consumes a parenthesis token
func (l *Lexer) consumeParenthesis() Token {
	t := Token{
		kind:   LeftParenthesis,
		column: l.column,
		line:   l.line,
		value:  "(",
	}

	if l.getCurr() == ')' {
		t.kind = RightParenthesis
		t.value = ")"
	}

	l.move()
	return t
}

// consumeOperator consumes an operator token
func (l *Lexer) consumeOperator() Token {
	char := l.getCurr()

	if isArithmeticOperator(char) {
		return l.consumeArithmeticOperator()
	}

	if isComparisonOperator(char) {
		return l.consumeComparisonOperator()
	}

	l.move()
	return UnknownToken
}

// consumeOperator consumes an operator token
func (l *Lexer) consumeArithmeticOperator() Token {
	defer l.move()
	t := Token{
		column: l.column,
		line:   l.line,
	}

	op := l.getCurr()

	switch op {
	case '+':
		t.kind = Plus
	case '-':
		t.kind = Minus
	case '/':
		t.kind = Div
	case '*':
		t.kind = Times
	default:
		return UnknownToken
	}

	t.value = string(op)

	return t
}

// consumeComparisonOperator consumes an operator token
func (l *Lexer) consumeComparisonOperator() Token {
	defer l.move()
	t := Token{
		column: l.column,
		line:   l.line,
	}

	char := l.getCurr()
	hasEquals := false

	if l.position+1 < len(l.input) {
		// copy next byte
		cpy := l.input[l.position+1]
		if cpy == '=' {
			hasEquals = true
			l.move()
		}
	}

	switch char {
	case '>':
		t.kind = LessThan
		if hasEquals {
			t.kind = LessThanOrEqual
		}
	case '<':
		t.kind = GreaterThan
		if hasEquals {
			t.kind = GreaterThanOrEqual
		}
	case '=':
		t.kind = Assign
		if hasEquals {
			t.kind = Equal
		}
	default:
		return UnknownToken
	}

	t.value = string(char)

	return t
}

// consumeIdentifier consumes an identifier and returns a token
func (l *Lexer) consumeIdentifier() Token {
	t := Token{
		kind:   Identifier,
		column: l.column,
		line:   l.line,
	}
	var identifier bytes.Buffer
	for l.position < len(l.input) {
		c := l.getCurr()
		if !isValidIdentifierChar(c) {
			break
		}

		identifier.WriteByte(c)
		l.move()
	}

	t.value = identifier.String()

	return t
}

/*
// TODO: assignent operator
// Other arithmetic operators
// Conditions
*/
