package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/amupitan/hero/ast"
	"github.com/amupitan/hero/ast/core"
	lx "github.com/amupitan/hero/lexer"
)

func expressionEqual(exp1, exp2 core.Expression) bool {
	// TODO(DEV) add Equals(expression) to core.Exp and use it here
	return exp1.String() == exp2.String()
}

func expectPanic(t *testing.T, want interface{}) {
	// strings are used here because a nil comparison
	// can't be used
	// see: https://golang.org/doc/faq#nil_error
	null, wantStr := fmt.Sprint(nil), fmt.Sprint(want)
	if r := recover(); wantStr == null && r != nil {
		// TODO(TEST) compare panic message
	} else if wantStr != null && r != nil {
		t.Errorf("Unexpected panic: %s", r.(error).Error())
	}
}

func TestParser_parse_expression(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  core.Expression
	}{
		{
			name:  "integer addition and multiplication",
			input: "1+2*3",
			want: &ast.Binary{
				Left:     &ast.Atom{Value: `1`, Type: lx.Int},
				Operator: lx.Token{Value: `+`, Type: lx.Plus},
				Right: &ast.Binary{
					Left:     &ast.Atom{Value: `2`, Type: lx.Int},
					Operator: lx.Token{Value: `*`, Type: lx.Times},
					Right:    &ast.Atom{Value: `3`, Type: lx.Int},
				},
			},
		},
		{
			name:  "integer addition and multiplication with parenthesis",
			input: "(1+2)*3",
			want: &ast.Binary{
				Left: &ast.Binary{
					Left:     &ast.Atom{Value: `1`, Type: lx.Int},
					Operator: lx.Token{Value: `+`, Type: lx.Plus},
					Right:    &ast.Atom{Value: `2`, Type: lx.Int},
				},
				Operator: lx.Token{Value: `*`, Type: lx.Times},
				Right:    &ast.Atom{Value: `3`, Type: lx.Int},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			if got := p.parse_expression(); !expressionEqual(got, tt.want) {
				t.Errorf("Parser.parse_expression() = %s,\n want %s", got, tt.want)
			}
		})
	}
}

func TestParser_parse_statement(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  core.Statement // if want is nil, a panic is expected
	}{
		{
			"variable declaration with type and value", // TODO(TEST) change test
			"var foo int = 0",
			&ast.Definition{Name: `foo`, Type: `int`, Value: &ast.Atom{Value: `0`, Type: lx.Int}},
		},
		{
			`variable declaration with no type or value`, // TODO(TEST) change test
			`var x`,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			defer expectPanic(t, tt.want)
			if got := p.parse_statement(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parse_statement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_attempt_parse_definition(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *ast.Definition
	}{
		{
			`variable declaration with type and value`,
			`var foo int = 0`,
			&ast.Definition{Name: `foo`, Type: `int`, Value: &ast.Atom{Value: `0`, Type: lx.Int}},
		},
		{
			`variable declaration with value`,
			`var bar = "hello"`,
			&ast.Definition{Name: `bar`, Value: &ast.Atom{Value: `hello`, Type: lx.String}},
		},
		{
			`variable declaration with type`,
			`var foobar int`,
			&ast.Definition{Name: `foobar`, Type: `int`},
		},
		{
			`variable declaration with no type or value`,
			`var x`,
			nil,
		},
		{
			`invalid declaration`,
			`var (invalid)`,
			nil,
		},
		{
			`short variable declaration with type and value`,
			`foo := 0`,
			&ast.Definition{Name: `foo`, Value: &ast.Atom{Value: `0`, Type: lx.Int}},
		},
		{
			`short variable declaration with invalid syntax`,
			`foo 0`,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			defer expectPanic(t, tt.want)
			if got := p.attempt_parse_definition(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.attempt_parse_definition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parse_binary(t *testing.T) {
	type fields struct {
		input string
		curr  int
	}
	type args struct {
		left  core.Expression
		my_op *lx.TokenType
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   core.Expression
	}{
		{
			name:   "integer addition",
			fields: fields{"1+2", 1},
			args:   args{left: &ast.Atom{Value: `1`, Type: lx.Int}, my_op: nil},
			want: &ast.Binary{
				Left:     &ast.Atom{Value: `1`, Type: lx.Int},
				Operator: lx.Token{Value: `+`, Type: lx.Plus, Line: 1, Column: 2},
				Right:    &ast.Atom{Value: `2`, Type: lx.Int},
			},
		},
		{
			name:   "integer addition and multiplication",
			fields: fields{"1+2*3", 1},
			args:   args{left: &ast.Atom{Value: `1`, Type: lx.Int}, my_op: nil},
			want: &ast.Binary{
				Left:     &ast.Atom{Value: `1`, Type: lx.Int},
				Operator: lx.Token{Value: `+`, Type: lx.Plus},
				Right: &ast.Binary{
					Left:     &ast.Atom{Value: `2`, Type: lx.Int},
					Operator: lx.Token{Value: `*`, Type: lx.Plus},
					Right:    &ast.Atom{Value: `3`, Type: lx.Int},
				},
			},
		},
		{
			name:   "integer assignment",
			fields: fields{"a = 4", 1},
			args:   args{left: &ast.Atom{Value: `a`, Type: lx.Identifier}, my_op: nil},
			want: &ast.Assignment{
				Identifier: `a`,
				Value:      &ast.Atom{Value: `4`, Type: lx.Int},
			},
		},
		{
			name:   "invalid double assignment",
			fields: fields{"a = 2 = 7", 1},
			args:   args{left: &ast.Atom{Value: `a`, Type: lx.Identifier}, my_op: nil},
			want:   nil,
		},
		{
			name:   "addition assignment",
			fields: fields{"a = 2 + .7", 1},
			args:   args{left: &ast.Atom{Value: `a`, Type: lx.Identifier}, my_op: nil},
			want: &ast.Assignment{
				Identifier: `a`,
				Value: &ast.Binary{
					Left:     &ast.Atom{Value: `2`, Type: lx.Int},
					Operator: lx.Token{Value: `+`, Type: lx.Plus},
					Right:    &ast.Atom{Value: `.7`, Type: lx.Float},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.fields.input)
			p.curr = tt.fields.curr
			defer expectPanic(t, tt.want)
			if got := p.parse_binary(tt.args.left, tt.args.my_op); !expressionEqual(got, tt.want) {
				t.Errorf("Parser.parse_binary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_attempt_parse_call(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        *ast.Call
		shouldPanic bool
	}{
		{
			name:  `call with no args`,
			input: `print()`,
			want: &ast.Call{
				Name: `print`,
				Args: []core.Expression{},
			},
		},
		{
			name:  `call with one arg`,
			input: `print(1)`,
			want: &ast.Call{
				Name: `print`,
				Args: []core.Expression{&ast.Atom{Type: `int`, Value: `1`}},
			},
		},
		{
			name:  `call with two args`,
			input: `print(1, "hello")`,
			want: &ast.Call{
				Name: `print`,
				Args: []core.Expression{&ast.Atom{Type: `int`, Value: `1`}, &ast.Atom{Type: `string`, Value: `hello`}},
			},
		},
		{
			name:        `call with extra separator at the end`,
			input:       `print(1, 2,)`,
			want:        nil,
			shouldPanic: true,
		},
		{
			name:        `call with extra separator at the end`,
			input:       `print(1, 2,)`,
			want:        nil,
			shouldPanic: true,
		},
		{
			name:  `invalid call - no parenthesis`,
			input: `print{1, "hello"}`,
			want:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			if tt.shouldPanic {
				defer expectPanic(t, nil)
			}
			if got := p.attempt_parse_call(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.attempt_parse_call() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_delimited(t *testing.T) {
	test_parser := func(p *Parser) core.Expression { return p.parse_expression() }
	type args struct {
		start       lx.TokenType
		stop        lx.TokenType
		separator   lx.TokenType
		end_sep     bool
		expr_parser parser
	}
	tests := []struct {
		name        string
		input       string
		args        args
		want        []core.Expression
		shouldPanic bool
	}{
		{
			name:  `empty braces`,
			input: `{}`,
			args:  args{lx.LeftBrace, lx.RightBrace, lx.Comma, false, test_parser},
			want:  []core.Expression{},
		},
		{
			name:  `one arg in bracket`,
			input: `[1]`,
			args:  args{lx.LeftBracket, lx.RightBracket, lx.Comma, false, test_parser},
			want:  []core.Expression{&ast.Atom{Type: `int`, Value: `1`}},
		},
		{
			name:  `two args in parenthesis`,
			input: `(1, "boom!")`,
			args:  args{lx.LeftParenthesis, lx.RightParenthesis, lx.Comma, false, test_parser},
			want:  []core.Expression{&ast.Atom{Type: `int`, Value: `1`}, &ast.Atom{Type: `string`, Value: `boom!`}},
		},
		{
			name:        `unexpected seperator at the end`,
			input:       `[1,]`,
			args:        args{lx.LeftBracket, lx.RightBracket, lx.Comma, false, test_parser},
			want:        nil,
			shouldPanic: true,
		},
		{
			name:  `expected seperator at the end`,
			input: `[1,]`,
			args:  args{lx.LeftBracket, lx.RightBracket, lx.Comma, true, test_parser},
			want:  []core.Expression{&ast.Atom{Type: `int`, Value: `1`}},
		},
		{
			name:  `numbers in braces`,
			input: `{1,2,3}`,
			args:  args{lx.LeftBrace, lx.RightBrace, lx.Comma, true, test_parser},
			want:  []core.Expression{&ast.Atom{Type: `int`, Value: `1`}, &ast.Atom{Type: `int`, Value: `2`}, &ast.Atom{Type: `int`, Value: `3`}},
		},
		{
			name:  `numbers in braces with separator at end`,
			input: `{1,2,3,}`,
			args:  args{lx.LeftBrace, lx.RightBrace, lx.Comma, true, test_parser},
			want:  []core.Expression{&ast.Atom{Type: `int`, Value: `1`}, &ast.Atom{Type: `int`, Value: `2`}, &ast.Atom{Type: `int`, Value: `3`}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			if tt.shouldPanic {
				defer expectPanic(t, nil)
			}
			if got := p.delimited(tt.args.start, tt.args.stop, tt.args.separator, tt.args.end_sep, tt.args.expr_parser); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.delimited() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parse_atom(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        core.Expression
		shouldPanic bool
	}{
		{
			name:  `integer`,
			input: `1`,
			want:  &ast.Atom{Value: `1`, Type: lx.Int},
		},
		{
			name:  `exponent float`,
			input: `1.2E-7`,
			want:  &ast.Atom{Value: `1.2E-7`, Type: lx.Float},
		},
		{
			name:  `raw string`,
			input: "`hello world`",
			want:  &ast.Atom{Value: `hello world`, Type: lx.RawString},
		},
		{
			name:  `rune`,
			input: `'爱'`,
			want:  &ast.Atom{Value: `爱`, Type: lx.Rune},
		},
		{
			name:  `underscore`,
			input: `_`,
			want:  &ast.Atom{Value: `_`, Type: lx.Underscore},
		},
		{
			name:  `single identifier`,
			input: `foo`,
			want:  &ast.Atom{Value: `foo`, Type: lx.Identifier},
		},
		{
			name:  `call object method with two args`,
			input: `foo.print(1, "hello")`,
			want: &ast.Call{
				Name:   `print`,
				Args:   []core.Expression{&ast.Atom{Type: `int`, Value: `1`}, &ast.Atom{Type: `string`, Value: `hello`}},
				Object: `foo`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			if tt.shouldPanic {
				defer expectPanic(t, nil)
			}
			if got := p.parse_atom(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parse_atom() = %v, want %v", got, tt.want)
			}
		})
	}
}
