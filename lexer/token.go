package lexer

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
	Plus  TokenType = "plus"  // +
	Minus TokenType = "minus" // -
	Times TokenType = "times" // *
	Div   TokenType = "div"   // /

	/// Comparison operators
	GreaterThan        TokenType = "greater than"             // >
	GreaterThanOrEqual TokenType = "greater than or equal to" // >=
	LessThan           TokenType = "less than"                // <
	LessThanOrEqual    TokenType = "less than or equal to"    // <=
	Equal              TokenType = "equals"                   // ==

	/// Assignment operator
	Assign TokenType = "assign" // =

	/// Parenthesis
	LeftParenthesis  TokenType = "left parenthesis"  // (
	RightParenthesis TokenType = "right parenthesis" // )

	/// Special TokenTypes
	EndOfInput TokenType = "end of input"
	Unknown    TokenType = "unknown"
)

var (
	UnknownToken    = Token{Unknown, string(Unknown), -1, -1}
	EndOfInputToken = Token{EndOfInput, string(EndOfInput), -1, -1}
)
