package lexer

import (
	"fmt"
)

type TokenType string

type Token struct {
	Type         TokenType
	Value        string
	Line, Column int
}

const (
	/// Identifiers and literals
	Identifier TokenType = "identifier"
	Int        TokenType = "int"
	Float      TokenType = "float"
	String     TokenType = "string"
	RawString  TokenType = "rstring"
	Rune       TokenType = "rune"
	Underscore TokenType = "_"

	/// Builtin Types
	BoolT    TokenType = "bool type"
	FloatT   TokenType = "float type"
	GenericT TokenType = "generic type"
	IntT     TokenType = "int type"
	RuneT    TokenType = "rune type"
	StringT  TokenType = "string type"

	/// Keywords
	Break     TokenType = "break"
	Class     TokenType = "class"
	Const     TokenType = "const"
	Continue  TokenType = "continue"
	Else      TokenType = "else"
	For       TokenType = "for"
	False     TokenType = "false"
	Func      TokenType = "func"
	If        TokenType = "if"
	Interface TokenType = "interface"
	Import    TokenType = "import"
	New_      TokenType = "new"
	Null      TokenType = "null"
	Package   TokenType = "package"
	Return    TokenType = "return"
	True      TokenType = "true"
	This      TokenType = "this"
	Var       TokenType = "var"

	/// Arithmetic operators
	Plus  TokenType = "+"
	Minus TokenType = "-"
	Times TokenType = "*"
	Div   TokenType = "/"
	Mod   TokenType = "%"

	PlusEq  TokenType = "+="
	MinusEq TokenType = "-="
	TimesEq TokenType = "*="
	DivEq   TokenType = "/="
	ModEq   TokenType = "%="

	Increment TokenType = "++"
	Decrement TokenType = "--"

	/// Comparison operators
	GreaterThan        TokenType = ">"
	GreaterThanOrEqual TokenType = ">="
	LessThan           TokenType = "<"
	LessThanOrEqual    TokenType = "<="
	Equal              TokenType = "=="
	NotEqual           TokenType = "!="

	/// Boolean operators
	And TokenType = "&&"
	Or  TokenType = "||"
	Not TokenType = "!"

	/// Bitwise operators
	BitAnd TokenType = "&"
	BitOr  TokenType = "|"
	BitXor TokenType = "^"
	BitNot TokenType = "~"

	BitAndEq TokenType = "&="
	BitOrEq  TokenType = "|="
	BitXorEq TokenType = "^="

	BitLeftShift  TokenType = "<<"
	BitRightShift TokenType = ">>"

	/// Assignment operators
	Assign  TokenType = "="
	Declare TokenType = ":="

	/// Dot
	Dot      TokenType = "."
	TwoDots  TokenType = ".."
	Ellipsis TokenType = "..."

	/// NewLine
	NewLine TokenType = `\n`

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
	Break:     struct{}{},
	Class:     struct{}{},
	Const:     struct{}{},
	Continue:  struct{}{},
	Else:      struct{}{},
	For:       struct{}{},
	False:     struct{}{},
	Func:      struct{}{},
	If:        struct{}{},
	Interface: struct{}{},
	Import:    struct{}{},
	New_:      struct{}{},
	Null:      struct{}{},
	Package:   struct{}{},
	Return:    struct{}{},
	True:      struct{}{},
	This:      struct{}{},
	Var:       struct{}{},
}

var (
	UnknownToken    = func(value string, Line, Column int) Token { return Token{Unknown, value, Line, Column} }
	EndOfInputToken = Token{EndOfInput, string(EndOfInput), -1, -1}
)

func (t Token) String() string {
	return fmt.Sprintf("Token(Value: %s, Type: %s, Position: %d:%d)", t.Value, t.Type, t.Line, t.Column)
}

// IsKeyword returns true if the token type is a keyword
func IsKeyword(t TokenType) bool {
	if _, ok := keywords[t]; ok {
		return true
	}
	return false
}

// IsBuiltinType returns true of the token type represents
// a builtin type
func IsBuiltinType(t TokenType) bool {
	return t == BoolT || t == FloatT || t == GenericT || t == IntT || t == RuneT || t == StringT
}
