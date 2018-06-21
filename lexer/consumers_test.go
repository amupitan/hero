package lexer

import (
	"reflect"
	"testing"
)

func TestLexer_consumableBoolOperator(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Token
	}{
		{
			name:  `&& operator`,
			input: `&&`,
			want:  Token{Type: And, Line: 1, Column: 1, Value: string(And)},
		},
		{
			name:  `|| operator`,
			input: `||`,
			want:  Token{Type: Or, Line: 1, Column: 1, Value: string(Or)},
		},
		{
			name:  `! operator`,
			input: `!`,
			want:  Token{Type: Not, Line: 1, Column: 1, Value: string(Not)},
		},
		{
			name:  `!= operator`,
			input: `!=`,
			want:  Token{Type: NotEqual, Line: 1, Column: 1, Value: string(NotEqual)},
		},
		{
			name:  `! before token`,
			input: `!x`,
			want:  Token{Type: Not, Line: 1, Column: 1, Value: string(Not)},
		},
		{
			name:  `invalid token after &`,
			input: `&x`,
			want:  UnknownToken(`x`, 1, 2),
		},
		{
			name:  `invalid token after &`,
			input: `|x`,
			want:  UnknownToken(`x`, 1, 2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			if got := l.consumableBoolOperator(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lexer.consumableBoolOperator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_consumeColonOrDeclare(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Token
	}{
		{
			name:  `consume colon`,
			input: `:`,
			want:  Token{Type: Colon, Value: string(Colon), Line: 1, Column: 1},
		},
		{
			name:  `consume declare`,
			input: `:=`,
			want:  Token{Type: Declare, Value: string(Declare), Line: 1, Column: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			if got := l.consumeColonOrDeclare(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lexer.consumeColonOrDeclare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_consumeComparisonOperator(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		want         Token
		wantPosition int
	}{
		{
			name:         `less than`,
			input:        `<`,
			want:         Token{Type: LessThan, Value: string(LessThan), Line: 1, Column: 1},
			wantPosition: 1,
		},
		{
			name:         `less than or equal to`,
			input:        `<=`,
			want:         Token{Type: LessThanOrEqual, Value: string(LessThanOrEqual), Line: 1, Column: 1},
			wantPosition: 2,
		},
		{
			name:         `greater than`,
			input:        `>`,
			want:         Token{Type: GreaterThan, Value: string(GreaterThan), Line: 1, Column: 1},
			wantPosition: 1,
		},
		{
			name:         `greater than or equal to`,
			input:        `>=`,
			want:         Token{Type: GreaterThanOrEqual, Value: string(GreaterThanOrEqual), Line: 1, Column: 1},
			wantPosition: 2,
		},
		{
			name:         `equal to`,
			input:        `==`,
			want:         Token{Type: Equal, Value: string(Equal), Line: 1, Column: 1},
			wantPosition: 2,
		},
		{
			name:         `assignment operator`,
			input:        `=`,
			want:         Token{Type: Assign, Value: string(Assign), Line: 1, Column: 1},
			wantPosition: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			if got := l.consumeComparisonOperator(); !reflect.DeepEqual(got, tt.want) || l.position != tt.wantPosition {
				t.Errorf("Lexer.consumeComparisonOperator() = %v, want %v. Lexer cursor is %d but expected to be %d", got, tt.want, l.position, tt.wantPosition)
			}
		})
	}
}

func TestLexer_consumeParenthesis(t *testing.T) {
	type fields struct {
		input    []rune
		position int
		Line     int
		Column   int
	}
	tests := []struct {
		name   string
		fields fields
		want   Token
	}{
		{
			"Comsume left parenthesis",
			fields{[]rune("(hello)"), 0, 1, 1},
			Token{LeftParenthesis, "(", 1, 1},
		},
		{
			"Comsume right parenthesis",
			fields{[]rune("(hello)"), 6, 1, 6},
			Token{RightParenthesis, ")", 1, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				input:    tt.fields.input,
				position: tt.fields.position,
				Line:     tt.fields.Line,
				Column:   tt.fields.Column,
			}
			if got := l.consumeDelimeter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lexer.consumeParenthesis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_consumeRune(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		want         Token
		wantPosition int
	}{
		{
			name:         `ascii character`,
			input:        `'y'`,
			want:         Token{Type: Rune, Value: `y`, Line: 1, Column: 2},
			wantPosition: 3,
		},
		{
			name:         `unicode character (emoji)`,
			input:        `'😂'`,
			want:         Token{Type: Rune, Value: `😂`, Line: 1, Column: 2},
			wantPosition: 3,
		},
		{
			name:         `unicode character (non-english char)`,
			input:        `'爱'`,
			want:         Token{Type: Rune, Value: `爱`, Line: 1, Column: 2},
			wantPosition: 3,
		},
		{
			name:         `escape character`,
			input:        `'\n'`,
			want:         Token{Type: Rune, Value: `\n`, Line: 1, Column: 2},
			wantPosition: 4,
		},
		{
			name:         `no opening quote`,
			input:        `x'`,
			want:         UnknownToken(`x`, 1, 1),
			wantPosition: 1,
		},
		{
			name:         `no closing quote`,
			input:        `'x`,
			want:         UnknownToken(``, 1, 3),
			wantPosition: 3,
		},
		{
			name:         `single quote`,
			input:        `'`,
			want:         UnknownToken(``, 1, 2),
			wantPosition: 2,
		},
		{
			name:         `more than one rune in quote`,
			input:        `'xy'`,
			want:         UnknownToken(`y`, 1, 3),
			wantPosition: 3,
		},
		{
			name:         `no content`,
			input:        ``,
			want:         UnknownToken(``, 1, 1), //TODO(DEV) handle no content gracefully
			wantPosition: 1,
		},
		{
			name:         `empty rune`,
			input:        `''`,
			want:         UnknownToken(``, 1, 3), // TODO(DEV) should be a custom error
			wantPosition: 3,
		},
		// {
		// 	name:         `invalid escape character`, //TODO(TEST) not all escape sequences are valid
		// 	input:        `'\n'`,
		// 	want:         UnknownToken(`\c`, 1, 3),
		// 	wantPosition: 3,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			if got := l.consumeRune(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lexer.consumeRune() = %v, want %v", got, tt.want)
			}
			if l.position != tt.wantPosition {
				t.Errorf("Lexer.consumeRune() Lexer cursor is %d but expected to be %d", l.position, tt.wantPosition)
			}
		})
	}
}

func TestLexer_consumeArithmeticOrBitOperator(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		want         Token
		wantPosition int
	}{
		{
			name:         `not operator`,
			input:        `!`,
			want:         UnknownToken(`!`, 1, 1),
			wantPosition: 0,
		},
		{
			name:         `bit-wise or operator`,
			input:        `|`,
			want:         Token{Type: BitOr, Line: 1, Column: 1, Value: string(BitOr)},
			wantPosition: 1,
		},
		{
			name:         `bit-wise and-assign operator`,
			input:        `&=`,
			want:         Token{Type: BitAndEq, Line: 1, Column: 1, Value: string(BitAndEq)},
			wantPosition: 2,
		},
		{
			name:         `bit-wise xor-assign operator`,
			input:        `^=`,
			want:         Token{Type: BitXorEq, Line: 1, Column: 1, Value: string(BitXorEq)},
			wantPosition: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			if got := l.consumeArithmeticOrBitOperator(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lexer.consumeArithmeticOrBitOperator() = %v, want %v", got, tt.want)
			}
			if l.position != tt.wantPosition {
				t.Errorf("Lexer.consumeArithmeticOrBitOperator() Lexer cursor is %d but expected to be %d", l.position, tt.wantPosition)
			}
		})
	}
}
