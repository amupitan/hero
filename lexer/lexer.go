package lexer

import (
	"bytes"
	"fmt"

	"github.com/amupitan/hero/lexer/fsm"
)

// Lexer performs lexical analysis on an input
type Lexer struct {
	input                  []byte
	position, line, column int
}

func New(input string) *Lexer {
	return &Lexer{
		input:  []byte(input),
		line:   1,
		column: 1,
	}
}

/// NextToken returns the next recognized token or an error if none is found
func (l *Lexer) NextToken() Token {
	l.skipWhiteSpace()
	if l.position >= len(l.input) {
		return EndOfInputToken
	}

	curr := l.getCurr()

	if isParenthesis(curr) {
		return l.consumeParenthesis()
	}

	if beginsLiteral(curr) {
		return l.recognizeLiteral()
	}

	if isOperator(curr) {
		return l.consumeOperator()
	}

	if isDot(curr) {
		return l.consumeDot()
	}

	return UnknownToken(string(curr), l.line, l.column)
}

// Tokenize returns all the tokens or an error
func (l *Lexer) Tokenize() ([]Token, error) {
	var token Token
	tokens := []Token{}
	for token = l.NextToken(); token.kind != EndOfInput && token.kind != Unknown; {
		tokens = append(tokens, token)
		token = l.NextToken()
	}

	if token.kind == Unknown {
		return nil, fmt.Errorf("Unrecognized character '%s' on line %d, column %d.", token.value, token.line, token.column)
	}

	return tokens, nil
}

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

	if isLetter(b) {
		return l.consumeIdentifierOrKeyword()
	}

	if beginsIdentifier(b) {
		// TODO: use comsumeIdentifier function to optmize out unused keyword functionality
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

	return l.consumableIdentifier(word)
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
	return Token{
		kind:   Dot,
		value:  string(Dot),
		line:   l.line,
		column: l.column,
	}
}

// getNextWord reads all the chracters till the next white space
// and returns the consumed characters
func (l *Lexer) getNextWord(isAllowed func(b byte) bool) string {
	var word bytes.Buffer
	if isAllowed == nil {
		isAllowed = func(b byte) bool { return true }
	}

	var i int
	for i = l.position; i < len(l.input); i++ {
		b := l.input[i]
		if isWhitespace(b) || !isAllowed(b) { //TODO: only spaces & tabs should count as whitespace
			break
		}
		word.WriteByte(b)
	}

	return word.String()
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

	// ignores whether token is found because we can
	// guarantee that atleast the first one will be found
	// otherwise this never would have been called
	num, _ := fsm.Run(l.input[l.position:])

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

func (l *Lexer) getUnknownToken(value string) Token {
	return UnknownToken(value, l.line, l.column)
}

// peek returns the byte at cursor and true if found,
// else it returns 0 and false
func (l *Lexer) peek() (byte, bool) {
	if l.position < len(l.input) {
		return l.input[l.position], true
	}
	return 0, false
}

// skipWhiteSpace skips all white spaces and new lines till the next non-space byte
func (l *Lexer) skipWhiteSpace() {
	for c, ok := l.peek(); ok && isWhitespace(c); c, ok = l.peek() {
		l.position++
		l.column++
		if c == '\n' {
			l.line++
			l.column = 1
		}
	}
}

/*
// TODO: assignment operator
// Other arithmetic operators
// Conditions
*/
