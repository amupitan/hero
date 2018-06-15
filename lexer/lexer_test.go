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
		Line     int
		Column   int
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
				Line:     tt.fields.Line,
				Column:   tt.fields.Column,
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
		Line     int
		Column   int
	}
	tests := []struct {
		name   string
		fields fields
		want   *Lexer
	}{
		{
			"only position and Column are updated after calling move",
			fields{[]byte("test"), 1, 1, 1},
			&Lexer{[]byte("test"), 2, 1, 2},
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
				Line:     tt.fields.Line,
				Column:   tt.fields.Column,
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
		Line     int
		Column   int
	}
	tests := []struct {
		name   string
		fields fields
		want   *Lexer
	}{
		{
			"space between characters on first Line",
			fields{[]byte("a = 3"), 1, 1, 2},
			&Lexer{[]byte("a = 3"), 2, 1, 3},
		},
		{
			"space between characters on another Line",
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
				Line:     tt.fields.Line,
				Column:   tt.fields.Column,
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
				Token{Column: 1, Type: Int, Line: 1, Value: "1"},
				Token{Column: 3, Type: Plus, Line: 1, Value: "+"},
				Token{Column: 5, Type: Int, Line: 1, Value: "1"},
				EndOfInputToken,
			},
			nil,
		},
		{
			"two dots and int",
			fields{"..3"},
			[]Token{
				Token{Column: 1, Type: TwoDots, Line: 1, Value: ".."},
				Token{Column: 3, Type: Int, Line: 1, Value: "3"},
				EndOfInputToken,
			},
			nil,
		},
		{
			"int ending with two dots",
			fields{"3.."},
			[]Token{
				Token{Column: 1, Type: Float, Line: 1, Value: "3."},
				Token{Column: 3, Type: Dot, Line: 1, Value: "."},
				EndOfInputToken,
			},
			nil,
		},
		{
			"float ending with two dots",
			fields{"3..."},
			[]Token{
				Token{Column: 1, Type: Float, Line: 1, Value: "3."},
				Token{Column: 3, Type: TwoDots, Line: 1, Value: ".."},
				EndOfInputToken,
			},
			nil,
		},
		{
			"float starting with a dot and dot",
			fields{".3."},
			[]Token{
				Token{Column: 1, Type: Float, Line: 1, Value: ".3"},
				Token{Column: 3, Type: Dot, Line: 1, Value: "."},
				EndOfInputToken,
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
				Token{Column: 1, Type: Identifier, Line: 1, Value: "a"},
				Token{Column: 3, Type: Plus, Line: 1, Value: "+"},
				Token{Column: 5, Type: RawString, Line: 1, Value: `hello`},
				EndOfInputToken,
			},
			nil,
		},
		{
			"identifier-string addition",
			fields{`a + "hello"`},
			[]Token{
				Token{Column: 1, Type: Identifier, Line: 1, Value: "a"},
				Token{Column: 3, Type: Plus, Line: 1, Value: "+"},
				Token{Column: 5, Type: String, Line: 1, Value: `hello`},
				EndOfInputToken,
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
				Token{Column: 1, Type: Identifier, Line: 1, Value: "a"},
				Token{Column: 2, Type: TwoDots, Line: 1, Value: ".."},
				Token{Column: 4, Type: Identifier, Line: 1, Value: `value`},
				Token{Column: 10, Type: Assign, Line: 1, Value: `=`},
				Token{Column: 12, Type: Int, Line: 1, Value: `3`},
				EndOfInputToken,
			},
			nil,
		},
		{
			"identifier and identifier",
			fields{"a && b"},
			[]Token{
				Token{Column: 1, Type: Identifier, Line: 1, Value: "a"},
				Token{Column: 3, Type: And, Line: 1, Value: "&&"},
				Token{Column: 6, Type: Identifier, Line: 1, Value: "b"},
				EndOfInputToken,
			},
			nil,
		},
		{
			"identifier bit-or identifier",
			fields{"a | b"},
			[]Token{
				Token{Column: 1, Type: Identifier, Line: 1, Value: "a"},
				Token{Column: 3, Type: BitOr, Line: 1, Value: "|"},
				Token{Column: 5, Type: Identifier, Line: 1, Value: "b"},
				EndOfInputToken,
			},
			nil,
		},
		{
			"identifier bit-or-equals identifier",
			fields{"a |= b"},
			[]Token{
				Token{Column: 1, Type: Identifier, Line: 1, Value: "a"},
				Token{Column: 3, Type: BitOrEq, Line: 1, Value: "|="},
				Token{Column: 6, Type: Identifier, Line: 1, Value: "b"},
				EndOfInputToken,
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
				Token{Column: 1, Type: Identifier, Line: 1, Value: "a"},
				Token{Column: 3, Type: BitLeftShift, Line: 1, Value: "<<"},
				Token{Column: 6, Type: Identifier, Line: 1, Value: "b"},
				EndOfInputToken,
			},
			nil,
		},
		{
			"identifier right shift identifier",
			fields{`a >> b`},
			[]Token{
				Token{Column: 1, Type: Identifier, Line: 1, Value: "a"},
				Token{Column: 3, Type: BitRightShift, Line: 1, Value: ">>"},
				Token{Column: 6, Type: Identifier, Line: 1, Value: "b"},
				EndOfInputToken,
			},
			nil,
		},
		{
			"identifier post-increment",
			fields{`a++`},
			[]Token{
				Token{Column: 1, Type: Identifier, Line: 1, Value: "a"},
				Token{Column: 2, Type: Increment, Line: 1, Value: "++"},
				EndOfInputToken,
			},
			nil,
		},
		{
			"identifier pre-decrement",
			fields{`--a`},
			[]Token{
				Token{Column: 1, Type: Decrement, Line: 1, Value: "--"},
				Token{Column: 3, Type: Identifier, Line: 1, Value: "a"},
				EndOfInputToken,
			},
			nil,
		},
		{
			"identifier post-decrement with comment",
			fields{`a-- // decrement`},
			[]Token{
				Token{Column: 1, Type: Identifier, Line: 1, Value: "a"},
				Token{Column: 2, Type: Decrement, Line: 1, Value: "--"},
				EndOfInputToken,
			},
			nil,
		},
		{
			"identifier post-decrement with comment and new line",
			fields{"a *= .2 // decrement\n\treturn a"},
			[]Token{
				Token{Column: 1, Type: Identifier, Line: 1, Value: "a"},
				Token{Column: 3, Type: TimesEq, Line: 1, Value: "*="},
				Token{Column: 6, Type: Float, Line: 1, Value: ".2"},
				Token{Column: 21, Type: NewLine, Line: 1, Value: `\n`},
				Token{Column: 2, Type: Return, Line: 2, Value: "return"},
				Token{Column: 9, Type: Identifier, Line: 2, Value: "a"},
				EndOfInputToken,
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
		Line     int
		Column   int
	}
	tests := []struct {
		name   string
		fields fields
		want   *Lexer
	}{
		{
			"last Line with comment",
			fields{[]byte("a += 3 // this adds 3"), 7, 1, 8},
			&Lexer{[]byte("a += 3 // this adds 3"), 21, 1, 22},
		},
		{
			"Line with comment",
			fields{[]byte("a += 3 // this adds 3\nb = 3"), 7, 1, 8},
			&Lexer{[]byte("a += 3 // this adds 3\nb = 3"), 21, 1, 22},
		},
		{
			"Line with empty comment",
			fields{[]byte("a += 3 //\nb = 3"), 7, 1, 8},
			&Lexer{[]byte("a += 3 //\nb = 3"), 9, 1, 10},
		},
		{
			"space between characters on another Line",
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
				Line:     tt.fields.Line,
				Column:   tt.fields.Column,
			}
			if l.skipComments(); !reflect.DeepEqual(l, tt.want) {
				t.Errorf("l = %#v, want %#v", l, tt.want)
			}
		})
	}
}
