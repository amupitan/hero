package lexer

import (
	"bytes"

	"github.com/amupitan/hero/lexer/fsm"
)

// consumeDelimeter consumes a delimeter token
func (l *Lexer) consumeDelimeter() Token {
	c := l.getCurr()
	t := Token{
		column: l.column,
		line:   l.line,
		value:  string(c),
	}

	switch c {
	case ',':
		t.kind = Comma
	case '(':
		t.kind = LeftParenthesis
	case ')':
		t.kind = RightParenthesis
	case '[':
		t.kind = LeftBracket
	case ']':
		t.kind = RightBracket
	case '{':
		t.kind = LeftBrace
	case '}':
		t.kind = RightBrace
	default:
		return UnknownToken(t.value, l.line, l.column)
	}

	l.move()
	return t
}

// consumeColonOrDeclare consumes a colon or declare token
func (l *Lexer) consumeColonOrDeclare() Token {
	t := Token{
		kind:   Colon,
		value:  string(Colon),
		column: l.column,
		line:   l.line,
	}

	l.move()

	// check if it is a `:=`
	if next, _ := l.peek(); next == '=' {
		t.kind = Declare
		t.value = `:=`
		l.move()
	}

	return t
}

// consumeOperator consumes an operator token
func (l *Lexer) consumeOperator() Token {
	defer l.move()
	char := l.getCurr()

	if isArithmeticOperator(char) {
		return l.consumeArithmeticOperator()
	}

	// if it isn't arithmetic then it is comparison
	return l.consumeComparisonOperator()
}

// consumeOperator consumes an operator token
func (l *Lexer) consumeArithmeticOperator() Token {
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
	}

	t.value = string(op)

	return t
}

// consumeComparisonOperator consumes an operator token
func (l *Lexer) consumeComparisonOperator() Token {
	t := Token{
		column: l.column,
		line:   l.line,
	}

	char := l.getCurr()
	hasEquals := false

	if l.position+1 < len(l.input) {
		// copy next byte
		cpy := l.input[l.position+1]

		// move cursor to accommodate '='
		if cpy == '=' {
			hasEquals = true
			l.move()
		}
	}

	switch char {
	case '<':
		if hasEquals {
			t.kind = LessThanOrEqual
			t.value = "<="
		} else {
			t.kind = LessThan
			t.value = "<"
		}
	case '>':
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
	}

	return t
}

func (l *Lexer) recognizeLiteral() Token {
	b := l.getCurr()

	if beginsIdentifier(b) {
		return l.consumeIdentifierOrKeyword()
	}

	if beginsNumber(b) {
		return l.consumeNumber()
	}

	if beginsString(b) {
		return l.consumeString()
	}

	if beginsRune(b) {
		return l.consumeRune()
	}

	return UnknownToken(string(b), l.line, l.column)

}

// consumeIdentifierOrKeyword recognizes an identifier or a keyword
func (l *Lexer) consumeIdentifierOrKeyword() Token {
	word := l.getNextWord(isValidIdentifierChar)
	defer func() {
		l.position += len(word)
		l.column += len(word)
	}()

	if t := l.consumableKeyword(word); t.kind != Unknown {
		return t
	}

	return Token{
		kind:   Identifier,
		value:  word,
		column: l.column,
		line:   l.line,
	}
}

// consumableKeyword returns a keyword/unknown token which can be consumed
func (l *Lexer) consumableKeyword(word string) Token {
	col, line := l.column, l.line

	keyword := TokenType(word)
	if _, ok := keywords[keyword]; ok {
		return Token{
			kind:   keyword,
			value:  word,
			column: col,
			line:   line,
		}
	}

	return UnknownToken(word, line, col)
}

// consumeDot consumes a keyword token
func (l *Lexer) consumeDot() Token {
	t := Token{
		kind:   Dot,
		value:  string(Dot),
		line:   l.line,
		column: l.column,
	}

	l.move()
	return t
}

// consumeRune consumes a rune token
func (l *Lexer) consumeRune() Token {
	if b := l.getCurr(); b != '\'' {
		t := l.getUnknownToken(string(b))
		l.move()
		return t
	}

	var value bytes.Buffer

	l.move()
	c := l.getCurr()

	// consume escape character if one exists
	if c == '\\' {
		value.WriteByte('\\')
		l.move()
		c = l.getCurr()
		// TODO: check valid escapes
	}
	l.move()

	if b := l.getCurr(); b != '\'' {
		t := l.getUnknownToken(string(b))
		l.move()
		return t
	}

	value.WriteByte(c)

	t := Token{
		column: l.column,
		line:   l.line,
		kind:   Rune,
		value:  value.String(),
	}
	l.move()
	return t
}

func (l *Lexer) consumeString() Token { return Token{} }

// consumableIdentifier returns an identifier/unknown token which can be consumed
func (l *Lexer) consumableIdentifier(word string) Token {
	t := Token{
		kind:   Identifier,
		column: l.column,
		line:   l.line,
	}

	for _, c := range word {
		if !isValidIdentifierChar(byte(c)) {
			break
		}
	}

	t.value = word
	return t
}

// consumeNumber consumes a number and returns an int or Float token
func (l *Lexer) consumeNumber() Token {
	fsm := fsm.New(states, states[0], nextState)

	num, isNum := fsm.Run(l.input[l.position:])
	if !isNum && len(num) == 0 {
		return UnknownToken(string(l.getCurr()), l.line, l.column)
	}

	// check for a decimal to determine whether Int or Float
	var kind TokenType = Int
	for _, b := range num {
		if b == '.' || b == 'e' || b == 'E' {
			kind = Float
		}
	}

	t := Token{
		kind:   kind,
		column: l.column,
		line:   l.line,
		value:  string(num),
	}
	l.position += len(num)
	l.column += len(num)

	return t
}
