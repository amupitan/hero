package lexer

import (
	"testing"
)

func TestToken_String(t *testing.T) {
	type fields struct {
		Type   TokenType
		Value  string
		Line   int
		Column int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   `identifier`,
			fields: fields{Identifier, `foo`, 1, 5},
			want:   `Token(Value: foo, Type: identifier, Position: 1:5)`,
		},
		{
			name:   `integer`,
			fields: fields{Int, `3`, 2, 1},
			want:   `Token(Value: 3, Type: int, Position: 2:1)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := Token{
				Type:   tt.fields.Type,
				Value:  tt.fields.Value,
				Line:   tt.fields.Line,
				Column: tt.fields.Column,
			}
			if got := tok.String(); got != tt.want {
				t.Errorf("Token.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsKeyword(t *testing.T) {
	tests := []struct {
		name  string
		token TokenType
		want  bool
	}{
		{
			name:  `return keyword`,
			token: Return,
			want:  true,
		},
		{
			name:  `an identifier is not a keyword`,
			token: Identifier,
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsKeyword(tt.token); got != tt.want {
				t.Errorf("IsKeyword() = %v, want %v", got, tt.want)
			}
		})
	}
}
