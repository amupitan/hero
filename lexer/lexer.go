package lexer

import (
	"bytes"
	"unicode"

	"./fsm"
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

func isLetter(b byte) bool              { return unicode.IsLetter(rune(b)) }
func isDigit(b byte) bool               { return unicode.IsDigit(rune(b)) }
func isValidIdentifierChar(b byte) bool { return false }
func isOperator(b byte) bool            { return false }
func isParenthesis(b byte) bool         { return b == '(' || b == ')' }
func isArithmeticOperator(b byte) bool  { return b == '+' || b == '-' || b == '*' || b == '/' }
func isComparisonOperator(b byte) bool  { return b == '>' || b == '<' || b == '=' }

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
		if hasEquals {
			t.kind = LessThanOrEqual
			t.value = "<="
		} else {
			t.kind = LessThan
			t.value = "<"
		}
	case '<':
		if hasEquals {
			t.kind = GreaterThanOrEqual
			t.value = ">="
		} else {
			t.kind = GreaterThan
			t.value = ">"
		}
	case '=':
		if hasEquals {
			t.kind = Equal
			t.value = "=="
		} else {
			t.kind = Assign
			t.value = "="
		}
	default:
		return UnknownToken
	}

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

func (l *Lexer) consumeFloat() Token {
	t := Token{
		kind:   Float,
		column: l.column,
		line:   l.line,
	}

	// fsm := fsm.New(states, states[0], nextState)
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

var (
	InitialState        = fsm.State{1, false}
	IntegerState        = fsm.State{2, true}
	BeginsFloatState    = fsm.State{3, false}
	FloatState          = fsm.State{4, true}
	BeginExpState       = fsm.State{5, false}
	BeginSignedExpState = fsm.State{6, false}
	ExponentState       = fsm.State{8, true}
	NullState           = fsm.NullState
)

var states = []fsm.State{
	InitialState,
	IntegerState,
	BeginsFloatState,
	FloatState,
	BeginExpState,
	BeginSignedExpState,
	ExponentState,
	NullState,
}

func nextState(currentState fsm.State, input byte) fsm.State {
	switch currentState.Value {
	case InitialState.Value:
		if isDigit(input) {
			return IntegerState
		}
	case IntegerState.Value:
		if isDigit(input) {
			return IntegerState
		}
		if input == '.' {
			return BeginsFloatState
		}
		if unicode.ToLower(rune(input)) == 'e' {
			return BeginExpState
		}
	case BeginsFloatState.Value:
		if isDigit(input) {
			return FloatState
		}
	case FloatState.Value:
		if isDigit(input) {
			return FloatState
		}
		if unicode.ToLower(rune(input)) == 'e' {
			return BeginExpState
		}
	case BeginExpState.Value:
		if isDigit(input) {
			return ExponentState
		}
		if input == '+' || input == '-' {
			return BeginSignedExpState
		}
	case BeginSignedExpState.Value:
		if isDigit(input) {
			return ExponentState
		}
	case ExponentState.Value:
		if isDigit(input) {
			return ExponentState
		}
	}
	return NullState
}

/*
// TODO: assignent operator
// Other arithmetic operators
// Conditions
*/
