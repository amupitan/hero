package lexer

import (
	"errors"
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
			BeginsFloatState,
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
			FloatState,
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
			if got := nextNumberState(tt.args.currentState, tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nextNumberState() = %v, want %v", got, tt.want)
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
			if got := l.consumeDelimeter(); !reflect.DeepEqual(got, tt.want) {
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
			&Lexer{[]byte("a = 3\nw * 3"), 5, 1, 5},
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
				t.Errorf("l = %#v, want %#v", l, tt.want)
			}
		})
	}
}

func TestLexer_Tokenize(t *testing.T) {
	type fields struct {
		input string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []Token
		wantErr error
	}{
		{
			"integer addition operation",
			fields{"1 + 1"},
			[]Token{
				Token{column: 1, kind: Int, line: 1, value: "1"},
				Token{column: 3, kind: Plus, line: 1, value: "+"},
				Token{column: 5, kind: Int, line: 1, value: "1"},
			},
			nil,
		},
		{
			"two dots and int",
			fields{"..3"},
			[]Token{
				Token{column: 1, kind: TwoDots, line: 1, value: ".."},
				Token{column: 3, kind: Int, line: 1, value: "3"},
			},
			nil,
		},
		{
			"int ending with two dots",
			fields{"3.."},
			[]Token{
				Token{column: 1, kind: Float, line: 1, value: "3."},
				Token{column: 3, kind: Dot, line: 1, value: "."},
			},
			nil,
		},
		{
			"float ending with two dots",
			fields{"3..."},
			[]Token{
				Token{column: 1, kind: Float, line: 1, value: "3."},
				Token{column: 3, kind: TwoDots, line: 1, value: ".."},
			},
			nil,
		},
		{
			"float starting with a dot and dot",
			fields{".3."},
			[]Token{
				Token{column: 1, kind: Float, line: 1, value: ".3"},
				Token{column: 3, kind: Dot, line: 1, value: "."},
			},
			nil,
		},
		{
			"symbol-integer addition operation",
			fields{"@ + 1"},
			nil,
			errors.New(`Unexpected token '@' on line 1, column 1.`),
		},
		{
			"identifier-raw_string addition",
			fields{"a + `hello`"},
			[]Token{
				Token{column: 1, kind: Identifier, line: 1, value: "a"},
				Token{column: 3, kind: Plus, line: 1, value: "+"},
				Token{column: 5, kind: RawString, line: 1, value: `hello`},
			},
			nil,
		},
		{
			"identifier-string addition",
			fields{`a + "hello"`},
			[]Token{
				Token{column: 1, kind: Identifier, line: 1, value: "a"},
				Token{column: 3, kind: Plus, line: 1, value: "+"},
				Token{column: 5, kind: String, line: 1, value: `hello`},
			},
			nil,
		},
		{
			"identifier-bad_string addition",
			fields{`a + "he"llo"`},
			nil,
			errors.New(`Unexpected token '"' on line 1, column 12.`),
		},
		{
			"identifier with double dots and assignment",
			fields{`a..value = 3`},
			[]Token{
				Token{column: 1, kind: Identifier, line: 1, value: "a"},
				Token{column: 2, kind: TwoDots, line: 1, value: ".."},
				Token{column: 4, kind: Identifier, line: 1, value: `value`},
				Token{column: 10, kind: Assign, line: 1, value: `=`},
				Token{column: 12, kind: Int, line: 1, value: `3`},
			},
			nil,
		},
		{
			"identifier and identifier",
			fields{"a && b"},
			[]Token{
				Token{column: 1, kind: Identifier, line: 1, value: "a"},
				Token{column: 3, kind: And, line: 1, value: "&&"},
				Token{column: 6, kind: Identifier, line: 1, value: "b"},
			},
			nil,
		},
		{
			"identifier bit-or identifier",
			fields{"a | b"},
			[]Token{
				Token{column: 1, kind: Identifier, line: 1, value: "a"},
				Token{column: 3, kind: BitOr, line: 1, value: "|"},
				Token{column: 5, kind: Identifier, line: 1, value: "b"},
			},
			nil,
		},
		{
			"identifier bit-or-equals identifier",
			fields{"a |= b"},
			[]Token{
				Token{column: 1, kind: Identifier, line: 1, value: "a"},
				Token{column: 3, kind: BitOrEq, line: 1, value: "|="},
				Token{column: 6, kind: Identifier, line: 1, value: "b"},
			},
			nil,
		},
		{
			"identifier not-equals identifier",
			fields{"a ~= b"},
			nil,
			errors.New(`Unexpected token '~' on line 1, column 3.`),
		},
		{
			"identifier left shift identifier",
			fields{`a << b`},
			[]Token{
				Token{column: 1, kind: Identifier, line: 1, value: "a"},
				Token{column: 3, kind: BitLeftShift, line: 1, value: "<<"},
				Token{column: 6, kind: Identifier, line: 1, value: "b"},
			},
			nil,
		},
		{
			"identifier right shift identifier",
			fields{`a >> b`},
			[]Token{
				Token{column: 1, kind: Identifier, line: 1, value: "a"},
				Token{column: 3, kind: BitRightShift, line: 1, value: ">>"},
				Token{column: 6, kind: Identifier, line: 1, value: "b"},
			},
			nil,
		},
		{
			"identifier post-increment",
			fields{`a++`},
			[]Token{
				Token{column: 1, kind: Identifier, line: 1, value: "a"},
				Token{column: 2, kind: Increment, line: 1, value: "++"},
			},
			nil,
		},
		{
			"identifier pre-decrement",
			fields{`--a`},
			[]Token{
				Token{column: 1, kind: Decrement, line: 1, value: "--"},
				Token{column: 3, kind: Identifier, line: 1, value: "a"},
			},
			nil,
		},
		{
			"identifier post-decrement with comment",
			fields{`a-- // decrement`},
			[]Token{
				Token{column: 1, kind: Identifier, line: 1, value: "a"},
				Token{column: 2, kind: Decrement, line: 1, value: "--"},
			},
			nil,
		},
		{
			"identifier post-decrement with comment and new line",
			fields{"a *= .2 // decrement\n\treturn a"},
			[]Token{
				Token{column: 1, kind: Identifier, line: 1, value: "a"},
				Token{column: 3, kind: TimesEq, line: 1, value: "*="},
				Token{column: 6, kind: Float, line: 1, value: ".2"},
				Token{column: 21, kind: NewLine, line: 1, value: `\n`},
				Token{column: 2, kind: Return, line: 2, value: "return"},
				Token{column: 9, kind: Identifier, line: 2, value: "a"},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.fields.input)
			got, err := l.Tokenize()
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Lexer.Tokenize() error = %s, wantErr %s", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lexer.Tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_skipComments(t *testing.T) {
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
			"last line with comment",
			fields{[]byte("a += 3 // this adds 3"), 7, 1, 8},
			&Lexer{[]byte("a += 3 // this adds 3"), 21, 1, 22},
		},
		{
			"line with comment",
			fields{[]byte("a += 3 // this adds 3\nb = 3"), 7, 1, 8},
			&Lexer{[]byte("a += 3 // this adds 3\nb = 3"), 21, 1, 22},
		},
		{
			"line with empty comment",
			fields{[]byte("a += 3 //\nb = 3"), 7, 1, 8},
			&Lexer{[]byte("a += 3 //\nb = 3"), 9, 1, 10},
		},
		{
			"space between characters on another line",
			fields{[]byte("a /= 3\nw * 3"), 2, 1, 3},
			&Lexer{[]byte("a /= 3\nw * 3"), 2, 1, 3},
		},
		{
			"expressions between lines",
			fields{[]byte("a = 3\nw * 3"), 5, 1, 6},
			&Lexer{[]byte("a = 3\nw * 3"), 5, 1, 6},
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
			if l.skipComments(); !reflect.DeepEqual(l, tt.want) {
				t.Errorf("l = %#v, want %#v", l, tt.want)
			}
		})
	}
}
