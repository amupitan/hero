package lexer

import (
	"unicode"

	"github.com/amupitan/hero/lexer/fsm"
)

var (
	InitialState        = fsm.State{1, false}
	IntegerState        = fsm.State{2, true}
	BeginsFloatState    = fsm.State{3, false}
	FloatState          = fsm.State{4, true}
	BeginExpState       = fsm.State{5, false}
	BeginSignedExpState = fsm.State{6, false}
	ExponentState       = fsm.State{8, true}
	NullState           = fsm.NullState
)

var states = []fsm.State{
	InitialState,
	IntegerState,
	BeginsFloatState,
	FloatState,
	BeginExpState,
	BeginSignedExpState,
	ExponentState,
	NullState,
}

func nextState(currentState fsm.State, input byte) fsm.State {
	switch currentState.Value {
	case InitialState.Value:
		if isDigit(input) {
			return IntegerState
		}
	case IntegerState.Value:
		if isDigit(input) {
			return IntegerState
		}
		if input == '.' {
			return BeginsFloatState
		}
		if unicode.ToLower(rune(input)) == 'e' {
			return BeginExpState
		}
	case BeginsFloatState.Value:
		if isDigit(input) {
			return FloatState
		}
	case FloatState.Value:
		if isDigit(input) {
			return FloatState
		}
		if unicode.ToLower(rune(input)) == 'e' {
			return BeginExpState
		}
	case BeginExpState.Value:
		if isDigit(input) {
			return ExponentState
		}
		if input == '+' || input == '-' {
			return BeginSignedExpState
		}
	case BeginSignedExpState.Value:
		if isDigit(input) {
			return ExponentState
		}
	case ExponentState.Value:
		if isDigit(input) {
			return ExponentState
		}
	}
	return NullState
}

func isLetter(b byte) bool              { return unicode.IsLetter(rune(b)) }
func isDigit(b byte) bool               { return unicode.IsDigit(rune(b)) }
func isValidIdentifierChar(b byte) bool { return b == '_' || isLetter(b) }
func isParenthesis(b byte) bool         { return b == '(' || b == ')' }
func isBitOrBoolOperator(b byte) bool   { return b == '&' || b == '|' || b == '!' }
func isComparisonOperator(b byte) bool  { return b == '>' || b == '<' || b == '=' }
func isWhitespace(b byte) bool          { return b == ' ' || b == '\t' || b == '\n' }
func beginsNumber(b byte) bool          { return b == '.' || isDigit(b) }
func beginsString(b byte) bool          { return b == '"' || b == '`' }
func beginsRune(b byte) bool            { return b == '\'' }
func beginsIdentifier(b byte) bool      { return b == '_' || isLetter(b) }

func beginsLiteral(b byte) bool {
	return beginsString(b) || beginsIdentifier(b) || beginsRune(b) || beginsNumber(b)
}

func isOperator(b byte) bool {
	return isArithmeticOperator(b) || isComparisonOperator(b) || isBitOrBoolOperator(b)
}
func isArithmeticOperator(b byte) bool {
	return b == '+' || b == '-' || b == '*' || b == '/' || b == '%'
}