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
			name:  "integer addition",
			input: "1+2",
			want: &ast.Binary{
				Left:     &ast.Atom{Value: `1`, Type: lx.Int},
				Operator: lx.Token{Value: `+`, Type: lx.Plus, Line: 1, Column: 2},
				Right:    &ast.Atom{Value: `2`, Type: lx.Int},
			},
		},
		{
			name:  "integer addition and multiplication",
			input: "1+2*3",
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
			name:  "integer addition and multiplication with parenthesis",
			input: "(1+2)*3",
			want: &ast.Binary{
				Left: &ast.Binary{
					Left:     &ast.Atom{Value: `1`, Type: lx.Int},
					Operator: lx.Token{Value: `+`, Type: lx.Plus},
					Right:    &ast.Atom{Value: `2`, Type: lx.Int},
				},
				Operator: lx.Token{Value: `*`, Type: lx.Plus},
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
