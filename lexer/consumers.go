package lexer

import (
	"bytes"

	"github.com/amupitan/hero/lexer/fsm"
)

// consumeDelimeter consumes a delimeter token
func (l *Lexer) consumeDelimeter() Token {
	c := l.getCurr()
	t := Token{
		Column: l.Column,
		Line:   l.Line,
		Value:  string(c),
	}

	switch c {
	case ',':
		t.Type = Comma
	case '(':
		t.Type = LeftParenthesis
	case ')':
		t.Type = RightParenthesis
	case '[':
		t.Type = LeftBracket
	case ']':
		t.Type = RightBracket
	case '{':
		t.Type = LeftBrace
	case '}':
		t.Type = RightBrace
	case ';':
		t.Type = SemiColon
	default:
		return UnknownToken(t.Value, l.Line, l.Column)
	}

	l.move()
	return t
}

// consumeNewline consumes a new line
func (l *Lexer) consumeNewline() Token {
	t := Token{
		Column: l.Column,
		Line:   l.Line,
		Type:   NewLine,
		Value:  `\n`,
	}

	l.position++
	l.Line++
	l.Column = 1

	return t
}

// consumeColonOrDeclare consumes a colon or declare token
func (l *Lexer) consumeColonOrDeclare() Token {
	t := Token{
		Type:   Colon,
		Value:  string(Colon),
		Column: l.Column,
		Line:   l.Line,
	}

	l.move()

	// check if it is a `:=`
	if next, _ := l.peek(); next == '=' {
		t.Type = Declare
		t.Value = `:=`
		l.move()
	}

	return t
}

// recognizeOperator consumes an operator token
func (l *Lexer) recognizeOperator() Token {
	c := l.getCurr()

	if isArithmeticOperator(c) || isBitOperator(c) || c == '!' {
		t := l.consumeArithmeticOrBitOperator()
		if t.Type == Unknown && isBoolOperator(c) {
			return l.consumableBoolOperator()
		}
		return t
	}

	// attempt to consume shift operator
	if beginsBitShift(c) {
		if t := l.consumeBitShiftOperator(); t.Type != Unknown {
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
		Column: l.Column,
		Line:   l.Line,
	}

	switch c {
	case '<':
		t.Type = BitLeftShift
		t.Value = string(BitLeftShift)
	case '>':
		t.Type = BitRightShift
		t.Value = string(BitRightShift)
	default:
		return l.getUnknownToken(string(c))
	}

	// consume first token
	l.move()
	// if the current and next tokens aren't the same
	// then it can't be a bit shift(<< or >>)
	if next, _ := l.peek(); c != next {
		t = UnknownToken(string(next), l.Line, l.Column)
		l.retract()
		return t
	}

	// consume second token
	l.move()
	return t
}

// consumeArithmeticOrBitOperator consumes an arithmetic or bit operator token
func (l *Lexer) consumeArithmeticOrBitOperator() Token {
	op := l.getCurr()
	t := Token{
		Column: l.Column,
		Line:   l.Line,
		Value:  string(op),
	}
	l.move()

	next, _ := l.peek()

	if next == '=' {
		switch op {
		case '+':
			t.Type = PlusEq
		case '-':
			t.Type = MinusEq
		case '/':
			t.Type = DivEq
		case '*':
			t.Type = TimesEq
		case '%':
			t.Type = ModEq
		case '&':
			t.Type = BitAndEq
		case '|':
			t.Type = BitOrEq
		case '^':
			t.Type = BitXorEq
		default:
			l.retract()
			return l.getUnknownToken(string(op))
		}

		// consume equals sign
		t.Value = string(op) + "="
		l.move()

		return t

	} else if !isBoolOperator(next) {
		switch op {
		case '+':
			t.Type = Plus
			// check if increment and consume
			if next == '+' {
				t.Type = Increment
				t.Value = "++"
				l.move()
			}
		case '-':
			t.Type = Minus
			// check if decrement and consume
			if next == '-' {
				t.Type = Decrement
				t.Value = "--"
				l.move()
			}
		case '/':
			t.Type = Div
		case '*':
			t.Type = Times
		case '%':
			t.Type = Mod
		case '&':
			t.Type = BitAnd
		case '|':
			t.Type = BitOr
		case '^':
			t.Type = BitXor
		case '~':
			t.Type = BitNot
		default:
			l.retract()
			return l.getUnknownToken(string(op))
		}
		return t
	}

	l.retract()
	return l.getUnknownToken(string(next))
}

// consumableBoolOperator consumes a bool operator token
func (l *Lexer) consumableBoolOperator() Token {
	t := Token{
		Column: l.Column,
		Line:   l.Line,
	}

	c := l.getCurr()
	l.move()
	next, _ := l.peek()

	if c != '!' && c != next {
		return l.getUnknownToken(string(next))
	}

	switch c {
	case '&':
		t.Type = And
		t.Value = string(And)
	case '|':
		t.Type = Or
		t.Value = string(Or)
	case '!':
		if next == '=' {
			t.Type = NotEqual
			t.Value = string(NotEqual)
		} else {
			t.Type = Not
			t.Value = string(Not)
		}
	}

	if t.Value != `!` {
		l.move()
	}
	return t
}

// consumeComparisonOperator consumes an operator token
func (l *Lexer) consumeComparisonOperator() Token {
	t := Token{
		Column: l.Column,
		Line:   l.Line,
	}

	char := l.getCurr()
	hasEquals := false

	if l.position+1 < len(l.input) {
		// copy next rune
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
			t.Type = LessThanOrEqual
			t.Value = "<="
		} else {
			t.Type = LessThan
			t.Value = "<"
		}
	case '>':
		if hasEquals {
			t.Type = GreaterThanOrEqual
			t.Value = ">="
		} else {
			t.Type = GreaterThan
			t.Value = ">"
		}
	case '=':
		if hasEquals {
			t.Type = Equal
			t.Value = "=="
		} else {
			t.Type = Assign
			t.Value = "="
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
		if t := l.consumeNumber(); t.Type != Unknown {
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

	return UnknownToken(string(b), l.Line, l.Column)

}

// consumeIdentifierOrKeyword recognizes an identifier or a keyword
func (l *Lexer) consumeIdentifierOrKeyword() Token {
	word := l.getNextWord(isValidIdentifierChar)
	defer func() {
		l.position += len(word)
		l.Column += len(word)
	}()

	if t := l.consumableKeyword(word); t.Type != Unknown {
		return t
	}

	Type := Identifier
	if word == `_` {
		Type = Underscore
	}

	return Token{
		Type:   Type,
		Value:  word,
		Column: l.Column,
		Line:   l.Line,
	}
}

// consumableKeyword returns a keyword/unknown token which can be consumed
// this also consumes true/false literals
func (l *Lexer) consumableKeyword(word string) Token {
	t := Token{
		Value:  word,
		Column: l.Column,
		Line:   l.Line,
	}

	keyword := TokenType(word)
	if keyword == `true` || keyword == `false` {
		t.Type = Bool
	} else if _, ok := keywords[keyword]; ok {
		t.Type = keyword
	} else {
		t.Type = Unknown
	}

	return t
}

// consumeDots consumes a dot or dots token
func (l *Lexer) consumeDots() Token {
	t := Token{
		Type:   Dot,
		Value:  string(Dot),
		Line:   l.Line,
		Column: l.Column,
	}
	l.move()

	// check for potential second dot to form two dots
	if next, _ := l.peek(); isDot(next) {
		t.Type = TwoDots
		t.Value = string(TwoDots)
		l.move()

		// check for potential third dot to form ellipsis
		if next, _ = l.peek(); isDot(next) {
			t.Type = Ellipsis
			t.Value = string(Ellipsis)
			l.move()
		}
	}

	return t
}

// consumeRune consumes a rune token
func (l *Lexer) consumeRune() Token {
	// consume_quote returns an empty Token and true if a quote
	// can be consumed else it returns an unknown token and false
	consume_quote := func() (Token, bool) {
		if b, ok := l.peek(); !ok || b != '\'' {
			col := l.Column
			l.move()
			if !ok {
				return UnknownToken(``, l.Line, col), false
			}
			return UnknownToken(string(b), l.Line, col), false
		}
		return Token{}, true
	}

	if t, ok := consume_quote(); !ok {
		return t
	}

	var value bytes.Buffer

	// consume opening quote
	l.move()

	// check character
	c, ok := l.peek()
	if !ok {
		col := l.Column
		l.move()
		return UnknownToken(``, l.Line, col)
	}

	col := l.Column
	// consume escape character if one exists
	if c == '\\' {
		value.WriteByte('\\')
		l.move()
		if c, ok = l.peek(); !ok {
			l.move()
			return l.getUnknownToken(``)
		}
		// TODO: check valid escapes
	}

	// write charcter
	value.WriteRune(c)

	// consume character
	l.move()

	if t, ok := consume_quote(); !ok {
		return t
	}

	// consume closing quote
	l.move()

	return Token{
		Column: col,
		Line:   l.Line,
		Type:   Rune,
		Value:  value.String(),
	}
}

func (l *Lexer) consumeString() Token {
	nextState := &nextStringState
	Type := String
	if l.getCurr() == '`' {
		nextState = &nextRawStringState
		Type = RawString
	}
	fsm := fsm.New(stringStates, stringStates[0], *nextState)

	buf, ok := fsm.Run(l.input[l.position:])
	if !ok {
		return UnknownToken(string(l.getCurr()), l.Line, l.Column)
	}

	length := buf.Len()

	// remove starting delimeter
	buf.ReadByte()
	// remove trailing delimeter
	buf.Truncate(length - 2)

	t := Token{
		Type:   Type,
		Column: l.Column,
		Line:   l.Line,
		Value:  buf.String(),
	}
	l.position += length
	l.Column += length

	return t
}

// consumableIdentifier returns an identifier/unknown token which can be consumed
func (l *Lexer) consumableIdentifier(word string) Token {
	t := Token{
		Type:   Identifier,
		Column: l.Column,
		Line:   l.Line,
	}

	for _, c := range word {
		if !isValidIdentifierChar(rune(c)) {
			break
		}
	}

	t.Value = word
	return t
}

// consumeNumber consumes a number and returns an int or Float token
func (l *Lexer) consumeNumber() Token {
	fsm := fsm.New(numberStates, numberStates[0], nextNumberState)

	buf, isNum := fsm.Run(l.input[l.position:])
	num := buf.String()
	if !isNum {
		return UnknownToken(string(l.getCurr()), l.Line, l.Column)
	}

	// check for a decimal/exponent to determine whether Int or Float
	var Type TokenType = Int
	for _, b := range num {
		if b == '.' || b == 'e' || b == 'E' {
			Type = Float
		}
	}

	t := Token{
		Type:   Type,
		Column: l.Column,
		Line:   l.Line,
		Value:  num,
	}
	l.position += len(num)
	l.Column += len(num)

	return t
}
