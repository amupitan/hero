package lexer

import (
	"bytes"
	"fmt"
)

// Lexer performs lexical analysis on an input
type Lexer struct {
	input                  []byte
	position, Line, Column int
}

const UnknownTokenError = `Unexpected token '%s' on line %d, column %d.`

func New(input string) *Lexer {
	return &Lexer{
		input:  []byte(input),
		Line:   1,
		Column: 1,
	}
}

/// NextToken returns the next recognized token or an error if none is found
func (l *Lexer) NextToken() Token {
	l.skipWhiteSpace()
	l.skipComments()
	if l.position >= len(l.input) {
		return EndOfInputToken
	}

	curr := l.getCurr()

	if isNewLine(curr) {
		return l.consumeNewline()
	}

	if isDelimeter(curr) {
		return l.consumeDelimeter()
	}

	if beginsLiteral(curr) {
		return l.recognizeLiteral()
	}

	if isColon(curr) {
		return l.consumeColonOrDeclare()
	}

	if isOperator(curr) {
		return l.recognizeOperator()
	}

	return UnknownToken(string(curr), l.Line, l.Column)
}

// Tokenize returns all the tokens or an error
func (l *Lexer) Tokenize() ([]Token, error) {
	var token Token
	tokens := []Token{}
	for token = l.NextToken(); token.Type != EndOfInput && token.Type != Unknown; token = l.NextToken() {
		tokens = append(tokens, token)
	}

	if token.Type == Unknown {
		return nil, fmt.Errorf(UnknownTokenError, token.Value, token.Line, token.Column)
	}

	return tokens, nil
}

// getCurr returns the byte at the current position
func (l *Lexer) getCurr() byte {
	return l.input[l.position]
}

// moves the cursor a step forward on the same line
func (l *Lexer) move() {
	l.position++
	l.Column++
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
	return UnknownToken(value, l.Line, l.Column)
}

// updateCursor adds offset to the position and column of the lexer's cursor
func (l *Lexer) updateCursor(offset int) {
	l.position += offset
	l.Column += offset
}

// retract moves the cursor one step back on the same line
func (l *Lexer) retract() {
	if l.position > 0 {
		l.position--
	}
	if l.Column > 0 {
		l.Column--
	}
}

// peek returns the byte at cursor and true if found,
// else it returns 0 and false
func (l *Lexer) peek() (byte, bool) {
	if l.position < len(l.input) {
		return l.input[l.position], true
	}
	return 0, false
}

// skipWhiteSpace skips all white spaces till the next non-space or newline byte
func (l *Lexer) skipWhiteSpace() {
	for c, ok := l.peek(); ok && isWhitespace(c); c, ok = l.peek() {
		l.position++
		l.Column++
	}
}

// skipComments skips any comments on the same line
func (l *Lexer) skipComments() {
	c, _ := l.peek()
	isComment := false

	// search for double slash
	if c == '/' {
		l.move()
		if c, _ = l.peek(); c == '/' {
			l.move()
			isComment = true
		} else {
			l.retract()
		}
	}

	// skip comment content
	if isComment {
		for c, ok := l.peek(); ok && !isNewLine(c); c, ok = l.peek() {
			// TODO(IMPROV) column increment doesn't have to be in the loop
			l.move()
		}
	}
}

/*
// hex, oct,unicode
// escape charcters
// Conditions
*/
