package lexer

type TokenType string

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
	GreaterThanOrEqual TokenType = "greater than or equal to" // >TokenType =
	LessThan           TokenType = "less than"                // <
	LessThanOrEqual    TokenType = "less than or equal to"    // <TokenType =
	Equal              TokenType = "equals"                   // TokenType =TokenType =

	/// Assignment operator
	Assign TokenType = "assign" // TokenType =

	/// Parenthesis
	LeftParenthesis  TokenType = "left parenthesis"  // (
	RightParenthesis TokenType = "right parenthesis" // )

	/// Special tokens
	EndOfInput TokenType = "end of input"
)
