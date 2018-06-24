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
				Body:        &ast.Block{},
			},
		},
		{
			name:  `3 args, multiple typesets, no return`,
			input: `func hello(x, y int, z MyType) {}`,
			want: &ast.Function{
				Definition:  ast.Definition{Name: `hello`, Type: string(lx.Func)},
				Parameters:  []*ast.Param{&ast.Param{Name: `x`, Type: types.Int}, &ast.Param{Name: `y`, Type: types.Int}, &ast.Param{Name: `z`, Type: CustomType(`MyType`)}},
				ReturnTypes: []types.Type{},
				Body:        &ast.Block{},
			},
		},
		{
			name:  `2 args, seaprate types, no return`,
			input: `func multiply(x int, y int) {}`,
			want: &ast.Function{
				Definition:  ast.Definition{Name: `multiply`, Type: string(lx.Func)},
				Parameters:  []*ast.Param{&ast.Param{Name: `x`, Type: types.Int}, &ast.Param{Name: `y`, Type: types.Int}},
				ReturnTypes: []types.Type{},
				Body:        &ast.Block{},
			},
		},
		{
			name:  `lambda - 2 args, joined type, no return`,
			input: `func(x, y int) {}`,
			want: &ast.Function{
				Definition:  ast.Definition{Type: string(lx.Func)},
				Parameters:  []*ast.Param{&ast.Param{Name: `x`, Type: types.Int}, &ast.Param{Name: `y`, Type: types.Int}},
				ReturnTypes: []types.Type{},
				Body:        &ast.Block{},
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
				Body:        &ast.Block{},
			},
		},
		{
			name:  `2 args, joined type, 2 returns`,
			input: `func compute(x, y int) (int, MyType) {}`,
			want: &ast.Function{
				Definition:  ast.Definition{Name: `compute`, Type: string(lx.Func)},
				Parameters:  []*ast.Param{&ast.Param{Name: `x`, Type: types.Int}, &ast.Param{Name: `y`, Type: types.Int}},
				ReturnTypes: []types.Type{types.Int, CustomType(`MyType`)},
				Body:        &ast.Block{},
			},
		},
		{
			name:  `1 arg, 1 statement, no return`,
			input: `func add2(x int) { a := x + 2 }`,
			want: &ast.Function{
				Definition:  ast.Definition{Name: `add2`, Type: string(lx.Func)},
				Parameters:  []*ast.Param{&ast.Param{Name: `x`, Type: types.Int}},
				ReturnTypes: []types.Type{},
				Body: &ast.Block{
					Statements: []core.Statement{&ast.Definition{
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
		},
		{
			name:  `no args, 1 return`,
			input: `func isCool() bool {}`,
			want: &ast.Function{
				Definition:  ast.Definition{Name: `isCool`, Type: string(lx.Func)},
				Parameters:  []*ast.Param{},
				ReturnTypes: []types.Type{types.Bool},
				Body:        &ast.Block{},
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
		{
			name:  `2 args, joined type, 1 return`,
			input: `func add(x, y int) int { return x + y }`,
			want: &ast.Function{
				Definition:  ast.Definition{Name: `add`, Type: string(lx.Func)},
				Parameters:  []*ast.Param{&ast.Param{Name: `x`, Type: types.Int}, &ast.Param{Name: `y`, Type: types.Int}},
				ReturnTypes: []types.Type{types.Int},
				Body: &ast.Block{
					Statements: []core.Statement{&ast.Return{
						Values: []core.Expression{&ast.Binary{
							Left:     &ast.Atom{Type: lx.Identifier, Value: `x`},
							Operator: lx.Token{Value: `+`, Type: lx.Plus, Line: 1, Column: 35},
							Right:    &ast.Atom{Type: lx.Identifier, Value: `y`},
						}},
					}},
				},
			},
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
		name      string
		input     string
		isNegated bool
		want      core.Expression
	}{
		{
			name:  `lambda declaration`,
			input: `func(x, y int) {}`,
			want: &ast.Function{
				Definition:  ast.Definition{Type: string(lx.Func)},
				Parameters:  []*ast.Param{&ast.Param{Name: `x`, Type: types.Int}, &ast.Param{Name: `y`, Type: types.Int}},
				ReturnTypes: []types.Type{},
				Body:        &ast.Block{},
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
				Body:        &ast.Block{},
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
					Body:        &ast.Block{},
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
					Body:        &ast.Block{},
					Lambda:      true,
				},
				Args: []core.Expression{&ast.Atom{Type: lx.Int, Value: `1`}, &ast.Atom{Type: lx.Identifier, Value: `z`}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			if got := p.attempt_parse_lambda_call(tt.isNegated); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.attempt_parse_lambda_call() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_attempt_parse_named_call(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		isNegated   bool
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
			if got := p.attempt_parse_named_call(tt.isNegated); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.attempt_parse_named_call() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parse_block(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *ast.Block
	}{
		{
			name:  `no statements`,
			input: `{}`,
			want:  &ast.Block{},
		},
		{
			name:  `1 statement`,
			input: `{x := y + 2}`,
			want: &ast.Block{
				Statements: []core.Statement{
					&ast.Definition{Name: `x`, Value: &ast.Binary{
						Left:     &ast.Atom{Value: `y`, Type: lx.Identifier},
						Operator: lx.Token{Value: `+`, Type: lx.Plus, Line: 1, Column: 9},
						Right:    &ast.Atom{Value: `2`, Type: lx.Int},
					}},
				},
			},
		},
		{
			name: `multiple statements`,
			input: `{
				var str string
				i := 2
				str = "ha" * i
			}`,
			want: &ast.Block{
				Statements: []core.Statement{
					&ast.Definition{
						Name: `str`,
						Type: string(lx.String),
					},
					&ast.Definition{
						Name:  `i`,
						Value: &ast.Atom{Value: `2`, Type: lx.Int},
					},
					&ast.Assignment{
						Identifier: `str`,
						Value: &ast.Binary{
							Left:     &ast.Atom{Value: `ha`, Type: lx.String},
							Operator: lx.Token{Value: `*`, Type: lx.Times, Line: 4, Column: 16},
							Right:    &ast.Atom{Value: `i`, Type: lx.Identifier},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			if got := p.parse_block(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parse_block() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parse_if(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        *ast.If
		shouldPanic bool
	}{
		{
			name:  `one variable evaluation, no parenthesis`,
			input: `if x {}`,
			want: &ast.If{
				Condition: &ast.Atom{Type: lx.Identifier, Value: `x`},
				Body:      &ast.Block{},
			},
		},
		{
			name:  `call evaluation, no parenthesis`,
			input: `if !isFree() {}`,
			want: &ast.If{
				Condition: &ast.Call{Name: `isFree`, Args: []core.Expression{}, Negated: true},
				Body:      &ast.Block{},
			},
		},
		{
			name:  `one variable evaluation with parenthesis`,
			input: `if (x) {}`,
			want: &ast.If{
				Condition: &ast.Atom{Type: lx.Identifier, Value: `x`},
				Body:      &ast.Block{},
			},
		},
		{
			name:  `multi-variable evaluation with parenthesis`,
			input: `if (x != getValue()) {}`,
			want: &ast.If{
				Condition: &ast.Binary{
					Left:     &ast.Atom{Type: lx.Identifier, Value: `x`},
					Operator: lx.Token{Value: `!=`, Type: lx.NotEqual, Line: 1, Column: 7},
					Right:    &ast.Call{Name: `getValue`, Args: []core.Expression{}},
				},
				Body: &ast.Block{},
			},
		},
		{
			name:  `multi-variable evaluation, no parenthesis`,
			input: `if x != getValue() {}`,
			want: &ast.If{
				Condition: &ast.Binary{
					Left:     &ast.Atom{Type: lx.Identifier, Value: `x`},
					Operator: lx.Token{Value: `!=`, Type: lx.NotEqual, Line: 1, Column: 6},
					Right:    &ast.Call{Name: `getValue`, Args: []core.Expression{}},
				},
				Body: &ast.Block{},
			},
		},
		{
			name:  `literal bool`,
			input: `if true {}`,
			want: &ast.If{
				Condition: &ast.Atom{Type: lx.Bool, Value: `true`},
				Body:      &ast.Block{},
			},
		},
		{
			name:        `invalid condition - integer`,
			input:       `if 3 {}`,
			shouldPanic: true,
		},
		{
			name:        `invalid condition - arithmetic operation`,
			input:       `if x + 3 {}`,
			shouldPanic: true,
		},
		{
			name:        `invalid condition - assignment`,
			input:       `if x = 3 {}`,
			shouldPanic: true,
		},
		{
			name: `negation evaluation, with two statements`,
			input: `
				if !x {
					x = y % 2
				}
			`,
			want: &ast.If{
				Condition: &ast.Atom{Type: lx.Identifier, Value: `x`, Negated: true},
				Body: &ast.Block{
					Statements: []core.Statement{&ast.Assignment{
						Identifier: `x`,
						Value: &ast.Binary{
							Left:     &ast.Atom{Type: lx.Identifier, Value: `y`},
							Operator: lx.Token{Value: `%`, Type: lx.Mod, Line: 3, Column: 12},
							Right:    &ast.Atom{Type: lx.Int, Value: `2`},
						},
					}},
				},
			},
		},
		{
			name: `if else`,
			input: `if x {} else{
				x = y
			}`,
			want: &ast.If{
				Condition: &ast.Atom{Type: lx.Identifier, Value: `x`},
				Body:      &ast.Block{},
				Else: &ast.If{
					Body: &ast.Block{
						Statements: []core.Statement{&ast.Assignment{Identifier: `x`, Value: &ast.Atom{Type: lx.Identifier, Value: `y`}}},
					},
				},
			},
		},
		{
			name:  `if else-if`,
			input: `if x {} else if (y) { x = y }`,
			want: &ast.If{
				Condition: &ast.Atom{Type: lx.Identifier, Value: `x`},
				Body:      &ast.Block{},
				Else: &ast.If{
					Condition: &ast.Atom{Type: lx.Identifier, Value: `y`},
					Body: &ast.Block{
						Statements: []core.Statement{&ast.Assignment{Identifier: `x`, Value: &ast.Atom{Type: lx.Identifier, Value: `y`}}},
					},
				},
			},
		},
		{
			name:  `if else-if else`,
			input: `if x {} else if (y) { x = y } else { runner.start() }`,
			want: &ast.If{
				Condition: &ast.Atom{Type: lx.Identifier, Value: `x`},
				Body:      &ast.Block{},
				Else: &ast.If{
					Condition: &ast.Atom{Type: lx.Identifier, Value: `y`},
					Body: &ast.Block{
						Statements: []core.Statement{&ast.Assignment{Identifier: `x`, Value: &ast.Atom{Type: lx.Identifier, Value: `y`}}},
					},
					Else: &ast.If{
						Body: &ast.Block{
							Statements: []core.Statement{&ast.Call{Name: `start`, Object: `runner`, Args: []core.Expression{}}},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			if tt.shouldPanic {
				defer expectPanic(t, nil)
			}
			if got := p.parse_if(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parse_if() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestParser_parse_return(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        *ast.Return
		shouldPanic bool
	}{
		{
			name:  `return nothing`,
			input: `return`,
			want:  &ast.Return{Values: []core.Expression{}},
		},
		{
			name:  `return one thing - a call`,
			input: `return getID()`,
			want: &ast.Return{Values: []core.Expression{
				&ast.Call{Name: `getID`, Args: []core.Expression{}},
			}},
		},
		{
			name:  `return values in parenthesis`,
			input: `return (1,2,3)`,
			want: &ast.Return{Values: []core.Expression{
				&ast.Atom{Type: lx.Int, Value: `1`},
				&ast.Atom{Type: lx.Int, Value: `2`},
				&ast.Atom{Type: lx.Int, Value: `3`},
			}},
		},
		{
			name:  `return values with no parenthesis`,
			input: `return 1,2,3`,
			want: &ast.Return{Values: []core.Expression{
				&ast.Atom{Type: lx.Int, Value: `1`},
				&ast.Atom{Type: lx.Int, Value: `2`},
				&ast.Atom{Type: lx.Int, Value: `3`},
			}},
		},
		{
			name:  `multi-type return`,
			input: `return isValid, "yea!", func(){}, 5.3e-9, true`, //TODO(TEST) add a bool to the returned things
			want: &ast.Return{Values: []core.Expression{
				&ast.Atom{Type: lx.Identifier, Value: `isValid`},
				&ast.Atom{Type: lx.String, Value: `yea!`},
				&ast.Function{
					Lambda:      true,
					Definition:  ast.Definition{Type: string(lx.Func)},
					Body:        &ast.Block{},
					Parameters:  []*ast.Param{},
					ReturnTypes: []types.Type{},
				},
				&ast.Atom{Type: lx.Float, Value: `5.3e-9`},
				&ast.Atom{Type: lx.Bool, Value: `true`},
			}},
		},
		{
			name:        `invalid token in return`,
			input:       `return 1, var`,
			shouldPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			if tt.shouldPanic {
				defer expectPanic(t, nil)
			}
			if got := p.parse_return(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parse_return() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parse_loop(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        ast.Loop
		shouldPanic bool
	}{
		{
			name:  `popular loop - implied types`,
			input: `for i := 0; i < s.length(); i++ {}`,
			want: &ast.ForLoop{
				PreLoop: &ast.Definition{
					Name:  `i`,
					Value: &ast.Atom{Type: lx.Int, Value: `0`},
				},
				Condition: &ast.Binary{
					Left:     &ast.Atom{Type: lx.Identifier, Value: `i`},
					Operator: lx.Token{Type: lx.LessThan, Value: `<`, Line: 1, Column: 15},
					Right: &ast.Call{
						Args:   []core.Expression{},
						Name:   `length`,
						Object: `s`,
					},
				},
				PostIteration: &ast.Assignment{
					Identifier: `i`,
					Value:      &ast.Operation{Type: lx.Increment},
				},
				Body: &ast.Block{},
			},
		},
		{
			name:  `empty for loop`,
			input: ` for {}`,
			want: &ast.ForLoop{
				Body: &ast.Block{},
			},
		},
		{
			name:  `empty for loop with semicolons`,
			input: ` for ;; {}`,
			want: &ast.ForLoop{
				Body: &ast.Block{},
			},
		},
		{
			name:  `for loop with only condition`,
			input: `for i == j {}`,
			want: &ast.ForLoop{
				Condition: &ast.Binary{
					Left:     &ast.Atom{Type: lx.Identifier, Value: `i`},
					Operator: lx.Token{Type: lx.Equal, Value: `==`, Line: 1, Column: 7},
					Right:    &ast.Atom{Type: lx.Identifier, Value: `j`},
				},
				Body: &ast.Block{},
			},
		},
		{
			name:  `for loop with only condition and semicolons`,
			input: `for ;i == j; {}`,
			want: &ast.ForLoop{
				Condition: &ast.Binary{
					Left:     &ast.Atom{Type: lx.Identifier, Value: `i`},
					Operator: lx.Token{Type: lx.Equal, Value: `==`, Line: 1, Column: 8},
					Right:    &ast.Atom{Type: lx.Identifier, Value: `j`},
				},
				Body: &ast.Block{},
			},
		},
		{
			name:  `for loop with initial statement and condition`,
			input: `for i = 0; i > 4; {}`,
			want: &ast.ForLoop{
				PreLoop: &ast.Assignment{Identifier: `i`, Value: &ast.Atom{Type: lx.Int, Value: `0`}},
				Condition: &ast.Binary{
					Left:     &ast.Atom{Type: lx.Identifier, Value: `i`},
					Operator: lx.Token{Type: lx.GreaterThan, Value: `>`, Line: 1, Column: 14},
					Right:    &ast.Atom{Type: lx.Int, Value: `4`},
				},
				Body: &ast.Block{},
			},
		},
		// {
		// 	name:       `for loop with initial statement and comma, and condition`,
		// 	input:      `for key, value := val; key < 5; {}`,
		// 	want:       (TODO), TODO(TEST) this should be a valid statement but multivaraibale assignment must be a thing first
		// },
		{
			name:        `for loop missing final semicolon`,
			input:       ` for i = 0; i < 4 {}`,
			shouldPanic: true,
		},
		{
			name:  `for loop with only initial expression`,
			input: ` for i = 0 ;; {}`,
			want: &ast.ForLoop{
				PreLoop: &ast.Assignment{Identifier: `i`, Value: &ast.Atom{Type: lx.Int, Value: `0`}},
				Body:    &ast.Block{},
			},
		},
		{
			name:  `for loop with only final expression`,
			input: ` for ;; i++ {}`,
			want: &ast.ForLoop{
				Body: &ast.Block{},
				PostIteration: &ast.Assignment{
					Identifier: `i`,
					Value:      &ast.Operation{Type: lx.Increment},
				},
			},
		},
		{
			name:  `empty for loop with body`,
			input: ` for { x-- }`,
			want: &ast.ForLoop{
				Body: &ast.Block{
					Statements: []core.Statement{
						&ast.Assignment{
							Identifier: `x`,
							Value:      &ast.Operation{Type: lx.Decrement},
						},
					},
				},
			},
		},
		{
			name:  `for range loop`,
			input: `for i, elem in array {}`,
			want: &ast.RangeLoop{
				First:    `i`,
				Second:   `elem`,
				Iterable: `array`,
				Body:     &ast.Block{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			if tt.shouldPanic {
				defer expectPanic(t, nil)
			}
			if got := p.parse_loop(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parse_loop() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestParser_parse_loop_range(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		want       *ast.RangeLoop
		wantCursor int
	}{
		{
			name:  `for range - 2 vars`,
			input: `for i, elem in array {}`,
			want: &ast.RangeLoop{
				First:    `i`,
				Second:   `elem`,
				Iterable: `array`,
				Body:     &ast.Block{},
			},
		},
		{
			name:  `for range - 1 var`,
			input: `for i in array {}`,
			want: &ast.RangeLoop{
				First:    `i`,
				Iterable: `array`,
				Body:     &ast.Block{},
			},
		},
		{
			name:  `for range with body`,
			input: `for i in array { x++ }`,
			want: &ast.RangeLoop{
				First:    `i`,
				Iterable: `array`,
				Body: &ast.Block{
					Statements: []core.Statement{
						&ast.Assignment{
							Identifier: `x`,
							Value:      &ast.Operation{Type: lx.Increment},
						},
					},
				},
			},
		},
		{
			name:       `for range comma with no second identifier`,
			input:      `for i, in array {}`,
			want:       nil,
			wantCursor: 1,
		},
		{
			name:       `for range comma with no in`,
			input:      `for key, value := val; key < 5; {}`,
			want:       nil,
			wantCursor: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			if got := p.attempt_parse_range_loop(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.attempt_parse_range_loop() = %v, want %v", got, tt.want)
			}
			if tt.want == nil && p.curr != tt.wantCursor {
				t.Errorf("Parser.attempt_parse_range_loop() cursor = %d, Expected cursor at %d", p.curr, tt.wantCursor)
			}
		})
	}
}
