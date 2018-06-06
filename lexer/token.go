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
	Rune       TokenType = "rune"
	Underscore TokenType = "_"

	/// Keywords
	Var   TokenType = "var"
	Const TokenType = "const"
	Null  TokenType = "null"
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

	/// Boolean operators
	And TokenType = "&&"
	Or  TokenType = "||"
	Not TokenType = "!"

	/// Bitwise operators
	BitAnd TokenType = "&"
	BitOr  TokenType = "|"
	BitXor TokenType = "^"

	/// Assignment operators
	Assign  TokenType = "="
	Declare TokenType = ":="

	/// Dot
	Dot TokenType = "."
	// TwoDots TokenType = ".."
	// Ellipses TokenType = "..."

	/// Delimeters
	Colon            TokenType = ":"
	Comma            TokenType = ","
	LeftParenthesis  TokenType = "("
	RightParenthesis TokenType = ")"
	LeftBracket      TokenType = "["
	RightBracket     TokenType = "]"
	LeftBrace        TokenType = "{"
	RightBrace       TokenType = "}"

	/// Special TokenTypes
	EndOfInput TokenType = "end of input"
	Unknown    TokenType = "unknown"
)

var keywords = map[TokenType]struct{}{
	Var:   struct{}{},
	Const: struct{}{},
	Null:  struct{}{},
	If:    struct{}{},
	While: struct{}{},
}

var (
	UnknownToken    = func(value string, line, column int) Token { return Token{Unknown, value, line, column} }
	EndOfInputToken = Token{EndOfInput, string(EndOfInput), -1, -1}
)

func (t Token) String() string {
	return fmt.Sprintf("Token(Value: %s, Type: %s, Position: %d:%d)", t.value, t.kind, t.line, t.column)
}
