package parser

import (
	"reflect"
	"testing"

	"github.com/amupitan/hero/ast"
	"github.com/amupitan/hero/ast/core"
	lx "github.com/amupitan/hero/lexer"
	"github.com/amupitan/hero/types"
)

func TestParser_parse_func(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        *ast.Function
		lambda      bool
		shouldPanic bool
	}{
		{
			name:  `2 args, joined type, no return`,
			input: `func add(x, y int) {}`,
			want: &ast.Function{
				Definition:  ast.Definition{Name: `add`, Type: string(lx.Func)},
				Parameters:  []*ast.Param{&ast.Param{Name: `x`, Type: types.Int}, &ast.Param{Name: `y`, Type: types.Int}},
				ReturnTypes: []types.Type{},
				Body:        []core.Statement{},
			},
		},
		{
			name:  `3 args, multiple typesets, no return`,
			input: `func hello(x, y int, z MyType) {}`,
			want: &ast.Function{
				Definition:  ast.Definition{Name: `hello`, Type: string(lx.Func)},
				Parameters:  []*ast.Param{&ast.Param{Name: `x`, Type: types.Int}, &ast.Param{Name: `y`, Type: types.Int}, &ast.Param{Name: `z`, Type: CustomType(`MyType`)}},
				ReturnTypes: []types.Type{},
				Body:        []core.Statement{},
			},
		},
		{
			name:  `2 args, seaprate types, no return`,
			input: `func multiply(x int, y int) {}`,
			want: &ast.Function{
				Definition:  ast.Definition{Name: `multiply`, Type: string(lx.Func)},
				Parameters:  []*ast.Param{&ast.Param{Name: `x`, Type: types.Int}, &ast.Param{Name: `y`, Type: types.Int}},
				ReturnTypes: []types.Type{},
				Body:        []core.Statement{},
			},
		},
		{
			name:  `lambda - 2 args, joined type, no return`,
			input: `func(x, y int) {}`,
			want: &ast.Function{
				Definition:  ast.Definition{Type: string(lx.Func)},
				Parameters:  []*ast.Param{&ast.Param{Name: `x`, Type: types.Int}, &ast.Param{Name: `y`, Type: types.Int}},
				ReturnTypes: []types.Type{},
				Body:        []core.Statement{},
				Lambda:      true,
			},
			lambda: true,
		},
		{
			name:  `2 args, joined type, 1 return`,
			input: `func equals(x, y int) bool {}`,
			want: &ast.Function{
				Definition:  ast.Definition{Name: `equals`, Type: string(lx.Func)},
				Parameters:  []*ast.Param{&ast.Param{Name: `x`, Type: types.Int}, &ast.Param{Name: `y`, Type: types.Int}},
				ReturnTypes: []types.Type{types.Bool},
				Body:        []core.Statement{},
			},
		},
		{
			name:  `2 args, joined type, 2 returns`,
			input: `func compute(x, y int) (int, MyType) {}`,
			want: &ast.Function{
				Definition:  ast.Definition{Name: `compute`, Type: string(lx.Func)},
				Parameters:  []*ast.Param{&ast.Param{Name: `x`, Type: types.Int}, &ast.Param{Name: `y`, Type: types.Int}},
				ReturnTypes: []types.Type{types.Int, CustomType(`MyType`)},
				Body:        []core.Statement{},
			},
		},
		{
			name:  `1 arg, 1 statement, no return`,
			input: `func add2(x int) { a := x + 2 }`,
			want: &ast.Function{
				Definition:  ast.Definition{Name: `add2`, Type: string(lx.Func)},
				Parameters:  []*ast.Param{&ast.Param{Name: `x`, Type: types.Int}},
				ReturnTypes: []types.Type{},
				Body: []core.Statement{&ast.Definition{
					Name: `a`,
					Value: &ast.Binary{
						Left:     &ast.Atom{Value: `x`, Type: lx.Identifier},
						Operator: lx.Token{Value: `+`, Type: lx.Plus, Line: 1, Column: 27},
						Right:    &ast.Atom{Value: `2`, Type: lx.Int},
					},
				},
				},
			},
		},
		{
			name:  `no args, 1 return`,
			input: `func isCool() bool {}`,
			want: &ast.Function{
				Definition:  ast.Definition{Name: `isCool`, Type: string(lx.Func)},
				Parameters:  []*ast.Param{},
				ReturnTypes: []types.Type{types.Bool},
				Body:        []core.Statement{},
			},
		},
		{
			name:        `no type with 1 arg`,
			input:       `func bad(x) bool {}`,
			shouldPanic: true,
		},
		{
			name:        `no type with args`,
			input:       `func bad(x, y) bool {}`,
			shouldPanic: true,
		},
		{
			name:        `extra seperator before type`,
			input:       `func bad(x, y, int) bool {}`,
			shouldPanic: true,
		},
		{
			name:        `extra seperator after type`,
			input:       `func bad(x, y int,) bool {}`,
			shouldPanic: true,
		},
		{
			name:        `only separator present`,
			input:       `func bad(,) bool {}`,
			shouldPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			if tt.shouldPanic {
				defer expectPanic(t, nil)
			}
			if got := p.parse_func(tt.lambda); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parse_func() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_attempt_parse_lambda_call(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  core.Expression
	}{
		{
			name:  `lambda declaration`,
			input: `func(x, y int) {}`,
			want: &ast.Function{
				Definition:  ast.Definition{Type: string(lx.Func)},
				Parameters:  []*ast.Param{&ast.Param{Name: `x`, Type: types.Int}, &ast.Param{Name: `y`, Type: types.Int}},
				ReturnTypes: []types.Type{},
				Body:        []core.Statement{},
				Lambda:      true,
			},
		},
		{
			name:  `lambda declaration with 2 returns`,
			input: `func(x, y SomeType) (MyType, bool) {}`,
			want: &ast.Function{
				Definition:  ast.Definition{Type: string(lx.Func)},
				Parameters:  []*ast.Param{&ast.Param{Name: `x`, Type: CustomType(`SomeType`)}, &ast.Param{Name: `y`, Type: CustomType(`SomeType`)}},
				ReturnTypes: []types.Type{CustomType(`MyType`), types.Bool},
				Body:        []core.Statement{},
				Lambda:      true,
			},
		},
		{
			name:  `lambda declaration call with no args`,
			input: `func() {}()`,
			want: &ast.Call{
				Func: &ast.Function{
					Definition:  ast.Definition{Type: string(lx.Func)},
					Parameters:  []*ast.Param{},
					ReturnTypes: []types.Type{},
					Body:        []core.Statement{},
					Lambda:      true,
				},
				Args: []core.Expression{},
			},
		},
		{
			name:  `lambda declaration call with 2 args`,
			input: `func(x, y int) {}(1, z)`,
			want: &ast.Call{
				Func: &ast.Function{
					Definition:  ast.Definition{Type: string(lx.Func)},
					Parameters:  []*ast.Param{&ast.Param{Name: `x`, Type: types.Int}, &ast.Param{Name: `y`, Type: types.Int}},
					ReturnTypes: []types.Type{},
					Body:        []core.Statement{},
					Lambda:      true,
				},
				Args: []core.Expression{&ast.Atom{Type: lx.Int, Value: `1`}, &ast.Atom{Type: lx.Identifier, Value: `z`}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			if got := p.attempt_parse_lambda_call(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.attempt_parse_lambda_call() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_attempt_parse_named_call(t *testing.T) {
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
			if got := p.attempt_parse_named_call(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.attempt_parse_named_call() = %v, want %v", got, tt.want)
			}
		})
	}
}
