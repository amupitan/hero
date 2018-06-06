package lexer

import (
	"bytes"
	"fmt"
)

// Lexer performs lexical analysis on an input
type Lexer struct {
	input                  []byte
	position, line, column int
}

const UnknownTokenError = `Unexpected token '%s' on line %d, column %d.`

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

	if isColon(curr) {
		return l.consumeColonOrDeclare()
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
		return nil, fmt.Errorf(UnknownTokenError, token.value, token.line, token.column)
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

func (l *Lexer) getUnknownToken(value string) Token {
	return UnknownToken(value, l.line, l.column)
}

// updateCursor adds offset to the position and column of the lexer's cursor
func (l *Lexer) updateCursor(offset int) {
	l.position += offset
	l.column += offset
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
