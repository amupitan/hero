package lexer

import (
	"fmt"
)

type TokenType string

type Token struct {
	kind         TokenType
	value        string
	line, column int
}

const (
	/// Identifiers and literals
	Identifier TokenType = "identifier"
	Int        TokenType = "int"
	Float      TokenType = "float"
	String     TokenType = "string"

	/// Conditions
	If    TokenType = "if"
	While TokenType = "while"

	/// Arithmetic operators
	Plus  TokenType = "+"
	Minus TokenType = "-"
	Times TokenType = "*"
	Div   TokenType = "/"
	Mod   TokenType = "%"

	/// Comparison operators
	GreaterThan        TokenType = ">"
	GreaterThanOrEqual TokenType = ">="
	LessThan           TokenType = "<"
	LessThanOrEqual    TokenType = "<="
	Equal              TokenType = "=="

	/// Assignment operator
	Assign TokenType = "="

	/// Parenthesis
	LeftParenthesis  TokenType = "("
	RightParenthesis TokenType = ")"

	/// Special TokenTypes
	EndOfInput TokenType = "end of input"
	Unknown    TokenType = "unknown"
)

var (
	UnknownToken    = Token{Unknown, string(Unknown), -1, -1}
	EndOfInputToken = Token{EndOfInput, string(EndOfInput), -1, -1}
)

func (t Token) String() string {
	return fmt.Sprintf("Token(Value: %s, Type: %s, Poistion: %d:%d)", t.value, t.kind, t.line, t.column)
}
