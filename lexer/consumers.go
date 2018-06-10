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
		t.Kind = Comma
	case '(':
		t.Kind = LeftParenthesis
	case ')':
		t.Kind = RightParenthesis
	case '[':
		t.Kind = LeftBracket
	case ']':
		t.Kind = RightBracket
	case '{':
		t.Kind = LeftBrace
	case '}':
		t.Kind = RightBrace
	default:
		return UnknownToken(t.value, l.line, l.column)
	}

	l.move()
	return t
}

// consumeNewline consumes a new line
func (l *Lexer) consumeNewline() Token {
	t := Token{
		column: l.column,
		line:   l.line,
		Kind:   NewLine,
		value:  `\n`,
	}

	l.position++
	l.line++
	l.column = 1

	return t
}

// consumeColonOrDeclare consumes a colon or declare token
func (l *Lexer) consumeColonOrDeclare() Token {
	t := Token{
		Kind:   Colon,
		value:  string(Colon),
		column: l.column,
		line:   l.line,
	}

	l.move()

	// check if it is a `:=`
	if next, _ := l.peek(); next == '=' {
		t.Kind = Declare
		t.value = `:=`
		l.move()
	}

	return t
}

// recognizeOperator consumes an operator token
func (l *Lexer) recognizeOperator() Token {
	c := l.getCurr()

	if isArithmeticOperator(c) || isBitOperator(c) || c == '!' {
		t := l.consumeArithmeticOrBitOperator()
		if t.Kind == Unknown && isBoolOperator(c) {
			return l.consumableBoolOperator()
		}
		return t
	}

	// attempt to consume shift operator
	if beginsBitShift(c) {
		if t := l.consumeBitShiftOperator(); t.Kind != Unknown {
			return t
		}
	}

	// if it isn't arithmetic, bit or boolean then it is comparison
	return l.consumeComparisonOperator()
}

// consumebitShiftOperator consumes a bit shifting operator
func (l *Lexer) consumeBitShiftOperator() Token {
	c := l.getCurr()
	t := Token{
		column: l.column,
		line:   l.line,
	}
	l.move()

	if next, _ := l.peek(); c != next {
		return l.getUnknownToken(string(next))
	}

	switch c {
	case '<':
		t.Kind = BitLeftShift
		t.value = string(BitLeftShift)
	case '>':
		t.Kind = BitRightShift
		t.value = string(BitRightShift)
	default:
		l.retract()
		return l.getUnknownToken(string(c))
	}

	l.move()
	return t
}

// consumeArithmeticOrBitOperator consumes an arithmetic or bit operator token
func (l *Lexer) consumeArithmeticOrBitOperator() Token {
	op := l.getCurr()
	t := Token{
		column: l.column,
		line:   l.line,
		value:  string(op),
	}
	l.move()

	next, _ := l.peek()

	if next == '=' {
		switch op {
		case '+':
			t.Kind = PlusEq
		case '-':
			t.Kind = MinusEq
		case '/':
			t.Kind = DivEq
		case '*':
			t.Kind = TimesEq
		case '%':
			t.Kind = ModEq
		case '&':
			t.Kind = BitAndEq
		case '|':
			t.Kind = BitOrEq
		case '^':
			t.Kind = BitXorEq
		default:
			l.retract()
			return l.getUnknownToken(string(op))
		}

		// consume equals sign
		t.value = string(op) + "="
		l.move()

		return t

	} else if !isBoolOperator(next) {
		switch op {
		case '+':
			t.Kind = Plus
			// check if increment and consume
			if next == '+' {
				t.Kind = Increment
				t.value = "++"
				l.move()
			}
		case '-':
			t.Kind = Minus
			// check if decrement and consume
			if next == '-' {
				t.Kind = Decrement
				t.value = "--"
				l.move()
			}
		case '/':
			t.Kind = Div
		case '*':
			t.Kind = Times
		case '%':
			t.Kind = Mod
		case '&':
			t.Kind = BitAnd
		case '|':
			t.Kind = BitOr
		case '^':
			t.Kind = BitXor
		case '~':
			t.Kind = BitNot
		}
		return t
	}

	l.retract()
	return l.getUnknownToken(string(next))
}

// consumableBoolOperator consumes a bool operator token
func (l *Lexer) consumableBoolOperator() Token {
	t := Token{
		column: l.column,
		line:   l.line,
	}

	c := l.getCurr()
	l.move()
	next, _ := l.peek()

	if c != '!' && c != next {
		return l.getUnknownToken(string(next))
	}

	switch c {
	case '&':
		t.Kind = And
		t.value = string(And)
	case '|':
		t.Kind = Or
		t.value = string(Or)
	case '!':
		t.Kind = Not
		t.value = string(Not)
	}

	if c != '!' {
		l.move()
	}
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
			t.Kind = LessThanOrEqual
			t.value = "<="
		} else {
			t.Kind = LessThan
			t.value = "<"
		}
	case '>':
		if hasEquals {
			t.Kind = GreaterThanOrEqual
			t.value = ">="
		} else {
			t.Kind = GreaterThan
			t.value = ">"
		}
	case '=':
		if hasEquals {
			t.Kind = Equal
			t.value = "=="
		} else {
			t.Kind = Assign
			t.value = "="
		}
	}

	l.move()
	return t
}

func (l *Lexer) recognizeLiteral() Token {
	b := l.getCurr()

	if beginsIdentifier(b) {
		return l.consumeIdentifierOrKeyword()
	}

	if beginsNumber(b) {
		if t := l.consumeNumber(); t.Kind != Unknown {
			return t
		}
		// if it began with a number literal, it is likely a dot
		return l.consumeDots()
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

	if t := l.consumableKeyword(word); t.Kind != Unknown {
		return t
	}

	return Token{
		Kind:   Identifier,
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
			Kind:   keyword,
			value:  word,
			column: col,
			line:   line,
		}
	}

	return UnknownToken(word, line, col)
}

// consumeDots consumes a dot or dots token
func (l *Lexer) consumeDots() Token {
	t := Token{
		Kind:   Dot,
		value:  string(Dot),
		line:   l.line,
		column: l.column,
	}
	l.move()

	// check for potential second dot to form two dots
	if next, _ := l.peek(); isDot(next) {
		t.Kind = TwoDots
		t.value = string(TwoDots)
		l.move()

		// check for potential third dot to form ellipsis
		if next, _ = l.peek(); isDot(next) {
			t.Kind = Ellipsis
			t.value = string(Ellipsis)
			l.move()
		}
	}

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
		Kind:   Rune,
		value:  value.String(),
	}
	l.move()
	return t
}

func (l *Lexer) consumeString() Token {
	nextState := &nextStringState
	kind := String
	if l.getCurr() == '`' {
		nextState = &nextRawStringState
		kind = RawString
	}
	fsm := fsm.New(stringStates, stringStates[0], *nextState)

	buf, ok := fsm.Run(l.input[l.position:])
	if !ok {
		return UnknownToken(string(l.getCurr()), l.line, l.column)
	}

	length := buf.Len()

	// remove starting delimeter
	buf.ReadByte()
	// remove trailing delimeter
	buf.Truncate(length - 2)

	t := Token{
		Kind:   kind,
		column: l.column,
		line:   l.line,
		value:  buf.String(),
	}
	l.position += length
	l.column += length

	return t
}

// consumableIdentifier returns an identifier/unknown token which can be consumed
func (l *Lexer) consumableIdentifier(word string) Token {
	t := Token{
		Kind:   Identifier,
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
	fsm := fsm.New(numberStates, numberStates[0], nextNumberState)

	buf, isNum := fsm.Run(l.input[l.position:])
	num := buf.String()
	if !isNum {
		return UnknownToken(string(l.getCurr()), l.line, l.column)
	}

	// check for a decimal/exponent to determine whether Int or Float
	var kind TokenType = Int
	for _, b := range num {
		if b == '.' || b == 'e' || b == 'E' {
			kind = Float
		}
	}

	t := Token{
		Kind:   kind,
		column: l.column,
		line:   l.line,
		value:  num,
	}
	l.position += len(num)
	l.column += len(num)

	return t
}
