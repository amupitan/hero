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
