package parser

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	lx "github.com/amupitan/hero/lexer"
)

// reports an error of reaching the end of the input will expecting
// another token
func (p *Parser) reportEndOfInput(expected *lx.TokenType) error {
	// TODO(DEV) add file name
	return fmt.Errorf("%d:%d: Expected `%s` but reached end of file.", p.Lexer.Line, p.Lexer.Column, *expected)
}

// reportUnexpected returns an error of receiving a wrong token type
func (p *Parser) reportUnexpected(expected *lx.TokenType) error {
	// TODO(DEV) add file name
	return fmt.Errorf("%d:%d: Expected `%v` but found `%s`.", p.Lexer.Line, p.Lexer.Column, *expected, p.peek().Value)
}

// reportUnexpectedMultiple returns an error for expecting one of a set of
// tokens
func (p *Parser) reportUnexpectedMultiple(expected ...lx.TokenType) error {
	// TODO(DEV) add file name
	sb := bytes.Buffer{}

	// Add line and column info
	sb.WriteString(strconv.Itoa(p.Lexer.Line)) // TODO(DEV) this or rune-int calculation
	sb.WriteRune(':')
	sb.WriteString(strconv.Itoa(p.Lexer.Column)) // TODO(DEV) this or rune-int calculation
	sb.WriteString(": Expected either ")
	for _, ex := range expected {
		sb.WriteString(string(ex))
		sb.WriteString(", ")
	}

	// remove last comma and space
	sb.Truncate(sb.Len() - 2)

	sb.WriteString(" but received ")
	if curr := p.peek(); curr != nil {
		sb.WriteString(curr.Value)
	} else {
		// TODO: should never get here
		sb.WriteString("Unknown")
	}

	return errors.New(sb.String())
}
