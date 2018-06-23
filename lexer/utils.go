package lexer

import (
	"unicode"

	"github.com/amupitan/hero/lexer/fsm"
)

var (
	// Number States
	InitialState        = fsm.State{1, false}
	IntegerState        = fsm.State{2, true}
	BeginsFloatState    = fsm.State{3, false}
	FloatState          = fsm.State{4, true}
	BeginExpState       = fsm.State{5, false}
	BeginSignedExpState = fsm.State{6, false}
	ExponentState       = fsm.State{8, true}

	// String States
	StringState    = fsm.State{9, false}
	EndStringState = fsm.State{10, true}

	// NullState
	NullState = fsm.NullState
)

var numberStates = []fsm.State{
	InitialState,
	IntegerState,
	BeginsFloatState,
	FloatState,
	BeginExpState,
	BeginSignedExpState,
	ExponentState,
	NullState,
}

var stringStates = []fsm.State{
	InitialState,
	StringState,
	EndStringState,
}

func nextNumberState(currentState fsm.State, input rune) fsm.State {
	switch currentState.Value {
	case InitialState.Value:
		if isDigit(input) {
			return IntegerState
		}
		if input == '.' {
			return BeginsFloatState
		}
	case IntegerState.Value:
		if isDigit(input) {
			return IntegerState
		}
		if input == '.' {
			return FloatState
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

func nextStringStateGenerator(delimeter rune) fsm.Transition {
	return func(currentState fsm.State, input rune) fsm.State {
		switch currentState {
		case InitialState:
			if input == delimeter {
				return StringState
			}
		case StringState:
			if input == delimeter {
				return EndStringState
			}
			return StringState
		}
		return NullState
	}
}

var nextStringState = nextStringStateGenerator('"')

var nextRawStringState = nextStringStateGenerator('`')

func isLetter(b rune) bool              { return unicode.IsLetter(rune(b)) }
func isDigit(b rune) bool               { return unicode.IsDigit(rune(b)) }
func isDot(b rune) bool                 { return b == '.' }
func isColon(b rune) bool               { return b == ':' }
func isValidIdentifierChar(b rune) bool { return b == '_' || isLetter(b) || isDigit(b) }
func isBoolOperator(b rune) bool        { return b == '&' || b == '|' || b == '!' }
func isComparisonOperator(b rune) bool  { return b == '>' || b == '<' || b == '=' }
func isWhitespace(b rune) bool          { return b == ' ' || b == '\t' }
func isNewLine(b rune) bool             { return b == '\n' }
func beginsNumber(b rune) bool          { return b == '.' || isDigit(b) }
func beginsBitShift(b rune) bool        { return b == '<' || b == '>' }
func beginsString(b rune) bool          { return b == '"' || b == '`' }
func beginsRune(b rune) bool            { return b == '\'' }
func beginsIdentifier(b rune) bool      { return b == '_' || isLetter(b) }

func isDelimeter(b rune) bool {
	return b == ',' || b == ';' || b == '(' || b == ')' || b == '[' || b == ']' || b == '{' || b == '}'
}

func beginsLiteral(b rune) bool {
	return beginsString(b) || beginsIdentifier(b) || beginsRune(b) || beginsNumber(b)
}

func isBitOperator(b rune) bool {
	return b == '&' || b == '|' || b == '~' || b == '^'
}

func isOperator(b rune) bool {
	return isArithmeticOperator(b) || isComparisonOperator(b) || isBoolOperator(b) || isBitOperator(b)
}
func isArithmeticOperator(b rune) bool {
	return b == '+' || b == '-' || b == '*' || b == '/' || b == '%'
}
