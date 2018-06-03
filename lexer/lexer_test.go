package lexer

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/amupitan/hero/lexer/fsm"
)

func TestNew(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want *Lexer
	}{
		{
			"create lexer",
			args{"this is code"},
			&Lexer{[]byte("this is code"), 0, 1, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nextState(t *testing.T) {
	type args struct {
		currentState fsm.State
		input        byte
	}
	tests := []struct {
		name string
		args args
		want fsm.State
	}{
		{
			"Initial with digit",
			args{InitialState, '8'},
			IntegerState,
		},
		{
			"Initial with decimal point",
			args{InitialState, '.'},
			NullState,
		},
		{
			"Initial with exponent,e",
			args{InitialState, 'e'},
			NullState,
		},
		{
			"Initial with exponent,E",
			args{InitialState, 'E'},
			NullState,
		},
		{
			"Initial with sign,+",
			args{InitialState, '+'},
			NullState,
		},
		{
			"Initial with sign,-",
			args{InitialState, '-'},
			NullState,
		},
		{
			"Integer with digit",
			args{IntegerState, '8'},
			IntegerState,
		},
		{
			"Integer with decimal point",
			args{IntegerState, '.'},
			BeginsFloatState,
		},
		{
			"Integer with exponent,e",
			args{IntegerState, 'e'},
			BeginExpState,
		},
		{
			"Integer with exponent,E",
			args{IntegerState, 'E'},
			BeginExpState,
		},
		{
			"Integer with sign,+",
			args{IntegerState, '+'},
			NullState,
		},
		{
			"Integer with sign,-",
			args{IntegerState, '-'},
			NullState,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextState(tt.args.currentState, tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nextState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_getCurr(t *testing.T) {
	type fields struct {
		input    []byte
		position int
		line     int
		column   int
	}
	input := []byte("1 + 2")
	randIdx := rand.Int() % len(input)
	tests := []struct {
		name   string
		fields fields
		want   byte
	}{
		{
			"get first byte",
			fields{input, 0, 0, 0},
			input[0],
		},
		{
			"get random byte",
			fields{input, randIdx, 0, 0},
			input[randIdx],
		},
		{
			"get last valid byte",
			fields{input, len(input) - 1, 0, 0},
			input[len(input)-1],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				input:    tt.fields.input,
				position: tt.fields.position,
				line:     tt.fields.line,
				column:   tt.fields.column,
			}
			if got := l.getCurr(); got != tt.want {
				t.Errorf("Lexer.getCurr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_move(t *testing.T) {
	type fields struct {
		input    []byte
		position int
		line     int
		column   int
	}
	tests := []struct {
		name   string
		fields fields
		want   *Lexer
	}{
		{
			"only position and column are updated after calling move",
			fields{[]byte("test"), 1, 1, 1},
			&Lexer{[]byte("test"), 2, 1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				input:    tt.fields.input,
				position: tt.fields.position,
				line:     tt.fields.line,
				column:   tt.fields.column,
			}
			l.move()
			if !reflect.DeepEqual(l, tt.want) {
				t.Errorf("l = %v, want %v", l, tt.want)
			}
		})
	}
}

func TestLexer_consumeParenthesis(t *testing.T) {
	type fields struct {
		input    []byte
		position int
		line     int
		column   int
	}
	tests := []struct {
		name   string
		fields fields
		want   Token
	}{
		{
			"Comsume left parenthesis",
			fields{[]byte("(hello)"), 0, 1, 1},
			Token{LeftParenthesis, "(", 1, 1},
		},
		{
			"Comsume right parenthesis",
			fields{[]byte("(hello)"), 6, 1, 6},
			Token{RightParenthesis, ")", 1, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				input:    tt.fields.input,
				position: tt.fields.position,
				line:     tt.fields.line,
				column:   tt.fields.column,
			}
			if got := l.consumeParenthesis(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lexer.consumeParenthesis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_skipWhiteSpace(t *testing.T) {
	type fields struct {
		input    []byte
		position int
		line     int
		column   int
	}
	tests := []struct {
		name   string
		fields fields
		want   *Lexer
	}{
		{
			"space between characters on first line",
			fields{[]byte("a = 3"), 1, 1, 2},
			&Lexer{[]byte("a = 3"), 2, 1, 3},
		},
		{
			"space between characters on another line",
			fields{[]byte("a = 3\nw * 3"), 7, 2, 2},
			&Lexer{[]byte("a = 3\nw * 3"), 8, 2, 3},
		},
		{
			"expressions between lines",
			fields{[]byte("a = 3\nw * 3"), 5, 1, 5},
			&Lexer{[]byte("a = 3\nw * 3"), 6, 2, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				input:    tt.fields.input,
				position: tt.fields.position,
				line:     tt.fields.line,
				column:   tt.fields.column,
			}
			l.skipWhiteSpace()
			if !reflect.DeepEqual(l, tt.want) {
				t.Errorf("l = %v, want %v", l, tt.want)
			}
		})
	}
}
