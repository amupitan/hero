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
		name   string
		input  string
		want   *ast.Function
		lambda bool
	}{
		{
			name:  `2 args, joined type, no return`,
			input: `func add(x, y int) {}`,
			want: &ast.Function{
				Definition:  ast.Definition{Name: `add`, Type: string(lx.Func)},
				Parameters:  []*ast.Param{&ast.Param{Name: `x`, Type: types.Int}, &ast.Param{Name: `y`, Type: types.Int}},
				ReturnTypes: []types.Type{},
				Body:        []core.Statement{},
				Owner:       nil,
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
				Owner:       nil,
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
				Owner:       nil,
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
				Owner:       nil,
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
				Owner:       nil,
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
				Owner:       nil,
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
				Owner: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			if got := p.parse_func(tt.lambda); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parse_func() = %v, want %v", got, tt.want)
			}
		})
	}
}
