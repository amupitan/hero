package parser

import (
	"reflect"
	"testing"

	"github.com/amupitan/hero/ast"
	"github.com/amupitan/hero/ast/core"
	lx "github.com/amupitan/hero/lexer"
	"github.com/amupitan/hero/types"
)

// expectPanic fails a test if there was no panic
// it is used as a defer
func expectPanic(t *testing.T, want interface{}) {
	if r := recover(); r != nil {
		// TODO(TEST) compare panic message with want
	} else {
		t.Errorf(`Expected a panic but there was no panic`)
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
				Operator: lx.Token{Value: `+`, Type: lx.Plus, Line: 1, Column: 2},
				Right: &ast.Binary{
					Left:     &ast.Atom{Value: `2`, Type: lx.Int},
					Operator: lx.Token{Value: `*`, Type: lx.Times, Line: 1, Column: 4},
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
					Operator: lx.Token{Value: `+`, Type: lx.Plus, Line: 1, Column: 3},
					Right:    &ast.Atom{Value: `2`, Type: lx.Int},
				},
				Operator: lx.Token{Value: `*`, Type: lx.Times, Line: 1, Column: 6},
				Right:    &ast.Atom{Value: `3`, Type: lx.Int},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			if got := p.parse_expression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parse_expression() = %s,\n want %s", got, tt.want)
			}
		})
	}
}

func TestParser_parse_statement(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        core.Statement // if want is nil, a panic is expected
		shouldPanic bool
	}{
		{
			`parse var`,
			`var foo int = 0`,
			&ast.Definition{Name: `foo`, Type: `int`, Value: &ast.Atom{Value: `0`, Type: lx.Int}},
			false,
		},
		{
			name:  `parse definition`,
			input: `x := y + 2`,
			want: &ast.Definition{Name: `x`, Value: &ast.Binary{
				Left:     &ast.Atom{Value: `y`, Type: lx.Identifier},
				Operator: lx.Token{Value: `+`, Type: lx.Plus, Line: 1, Column: 8},
				Right:    &ast.Atom{Value: `2`, Type: lx.Int},
			}},
		},
		{
			name:  "parse expression",
			input: "1+2*3",
			want: &ast.Binary{
				Left:     &ast.Atom{Value: `1`, Type: lx.Int},
				Operator: lx.Token{Value: `+`, Type: lx.Plus, Line: 1, Column: 2},
				Right: &ast.Binary{
					Left:     &ast.Atom{Value: `2`, Type: lx.Int},
					Operator: lx.Token{Value: `*`, Type: lx.Times, Line: 1, Column: 4},
					Right:    &ast.Atom{Value: `3`, Type: lx.Int},
				},
			},
		},
		{
			name:  `parse func`,
			input: `func compute(x, y int) (int, MyType) {}`,
			want: &ast.Function{
				Definition:  ast.Definition{Name: `compute`, Type: string(lx.Func)},
				Parameters:  []*ast.Param{&ast.Param{Name: `x`, Type: types.Int}, &ast.Param{Name: `y`, Type: types.Int}},
				ReturnTypes: []types.Type{types.Int, CustomType(`MyType`)},
				Body:        &ast.Block{},
			},
		},
		{
			name:  `parse block`,
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
			name:  `parse return`,
			input: `return 1,2`,
			want: &ast.Return{Values: []core.Expression{
				&ast.Atom{Type: lx.Int, Value: `1`},
				&ast.Atom{Type: lx.Int, Value: `2`},
			}},
		},
		{
			name: `parse if`,
			input: `if x >= 4 {
						return "yolo"
					}else{
						return "dope"
					}`,
			want: &ast.If{
				Condition: &ast.Binary{
					Left:     &ast.Atom{Type: lx.Identifier, Value: `x`},
					Operator: lx.Token{Value: `>=`, Type: lx.GreaterThanOrEqual, Line: 1, Column: 6},
					Right:    &ast.Atom{Type: lx.Int, Value: `4`},
				},
				Body: &ast.Block{
					Statements: []core.Statement{
						&ast.Return{
							Values: []core.Expression{&ast.Atom{Type: lx.String, Value: `yolo`}},
						},
					},
				},
				Else: &ast.If{
					Body: &ast.Block{
						Statements: []core.Statement{
							&ast.Return{
								Values: []core.Expression{&ast.Atom{Type: lx.String, Value: `dope`}},
							},
						},
					},
				},
			},
		},
		{
			name:  `parse for loop`,
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
			name: `parse named loop`,
			input: `flex:
			for i == j {}`,
			want: &ast.ForLoop{
				Name: `flex`,
				Condition: &ast.Binary{
					Left:     &ast.Atom{Type: lx.Identifier, Value: `i`},
					Operator: lx.Token{Type: lx.Equal, Value: `==`, Line: 2, Column: 10},
					Right:    &ast.Atom{Type: lx.Identifier, Value: `j`},
				},
				Body: &ast.Block{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			if tt.shouldPanic {
				defer expectPanic(t, nil)
			}
			if got := p.parse_statement(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parse_statement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_attempt_parse_definition(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        *ast.Definition
		shouldPanic bool
	}{
		{
			name:  `variable declaration with type and value`,
			input: `var foo int = 0`,
			want:  &ast.Definition{Name: `foo`, Type: `int`, Value: &ast.Atom{Value: `0`, Type: lx.Int}},
		},
		{
			name:  `variable declaration with value`,
			input: `var bar = "hello"`,
			want:  &ast.Definition{Name: `bar`, Value: &ast.Atom{Value: `hello`, Type: lx.String}},
		},
		{
			name:  `variable declaration with type`,
			input: `var foobar int`,
			want:  &ast.Definition{Name: `foobar`, Type: `int`},
		},
		{
			name:        `variable declaration with no type or value`,
			input:       `var x`,
			want:        nil,
			shouldPanic: true,
		},
		{
			name:        `invalid declaration`,
			input:       `var (invalid)`,
			want:        nil,
			shouldPanic: true,
		},
		{
			name:  `short variable declaration with type and value`,
			input: `foo := 0`,
			want:  &ast.Definition{Name: `foo`, Value: &ast.Atom{Value: `0`, Type: lx.Int}},
		},
		{
			name:  `short variable declaration to expression with type and value`,
			input: `x := y + 2`,
			want: &ast.Definition{Name: `x`, Value: &ast.Binary{
				Left:     &ast.Atom{Value: `y`, Type: lx.Identifier},
				Operator: lx.Token{Value: `+`, Type: lx.Plus, Line: 1, Column: 8},
				Right:    &ast.Atom{Value: `2`, Type: lx.Int},
			}},
		},
		{
			name:  `short variable declaration with invalid syntax`,
			input: `foo 0`,
			want:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			if tt.shouldPanic {
				defer expectPanic(t, nil)
			}
			if got := p.attempt_parse_definition(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.attempt_parse_definition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parse_binary(t *testing.T) {
	invalid_op := lx.TokenType(`#`)
	type fields struct {
		input string
		curr  int
	}
	type args struct {
		left  core.Expression
		my_op *lx.TokenType
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        core.Expression
		shouldPanic bool
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
				Operator: lx.Token{Value: `+`, Type: lx.Plus, Line: 1, Column: 2},
				Right: &ast.Binary{
					Left:     &ast.Atom{Value: `2`, Type: lx.Int},
					Operator: lx.Token{Value: `*`, Type: lx.Times, Line: 1, Column: 4},
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
			name:        "invalid double assignment",
			fields:      fields{"a = 2 = 7", 1},
			args:        args{left: &ast.Atom{Value: `a`, Type: lx.Identifier}, my_op: nil},
			shouldPanic: true,
		},
		{
			name:   "addition assignment",
			fields: fields{"a = 2 + .7", 1},
			args:   args{left: &ast.Atom{Value: `a`, Type: lx.Identifier}, my_op: nil},
			want: &ast.Assignment{
				Identifier: `a`,
				Value: &ast.Binary{
					Left:     &ast.Atom{Value: `2`, Type: lx.Int},
					Operator: lx.Token{Value: `+`, Type: lx.Plus, Line: 1, Column: 7},
					Right:    &ast.Atom{Value: `.7`, Type: lx.Float},
				},
			},
		},
		{
			name:   `invalid operator`,
			fields: fields{`1 + 2`, 1},
			args:   args{left: &ast.Atom{Value: `1`, Type: lx.Int}, my_op: &invalid_op},
			want:   &ast.Atom{Value: `1`, Type: lx.Int},
		},
		{
			name:   `increment operator`,
			fields: fields{`i++`, 1},
			args:   args{left: &ast.Atom{Value: `i`, Type: lx.Identifier}, my_op: nil},
			want: &ast.Assignment{
				Identifier: `i`,
				Value:      &ast.Operation{Type: lx.Increment},
			},
		},
		{
			name:   `decrement operator`,
			fields: fields{`i--`, 1},
			args:   args{left: &ast.Atom{Value: `i`, Type: lx.Identifier}, my_op: nil},
			want: &ast.Assignment{
				Identifier: `i`,
				Value:      &ast.Operation{Type: lx.Decrement},
			},
		},
		{
			name:   "integer plus-assignment",
			fields: fields{"a += 4", 1},
			args:   args{left: &ast.Atom{Value: `a`, Type: lx.Identifier}, my_op: nil},
			want: &ast.Assignment{
				Identifier: `a`,
				Value: &ast.Operation{
					Type:  lx.PlusEq,
					Value: &ast.Atom{Value: `4`, Type: lx.Int},
				},
			},
		},
		{
			name:   "integer times-assignment",
			fields: fields{"a *= 4", 1},
			args:   args{left: &ast.Atom{Value: `a`, Type: lx.Identifier}, my_op: nil},
			want: &ast.Assignment{
				Identifier: `a`,
				Value: &ast.Operation{
					Type:  lx.TimesEq,
					Value: &ast.Atom{Value: `4`, Type: lx.Int},
				},
			},
		},
		{
			name:   "integer div-assignment",
			fields: fields{"a /= 4", 1},
			args:   args{left: &ast.Atom{Value: `a`, Type: lx.Identifier}, my_op: nil},
			want: &ast.Assignment{
				Identifier: `a`,
				Value: &ast.Operation{
					Type:  lx.DivEq,
					Value: &ast.Atom{Value: `4`, Type: lx.Int},
				},
			},
		},
		{
			name:   "integer minus-assignment",
			fields: fields{"a -= 4", 1},
			args:   args{left: &ast.Atom{Value: `a`, Type: lx.Identifier}, my_op: nil},
			want: &ast.Assignment{
				Identifier: `a`,
				Value: &ast.Operation{
					Type:  lx.MinusEq,
					Value: &ast.Atom{Value: `4`, Type: lx.Int},
				},
			},
		},
		{
			name:   "integer mod-assignment",
			fields: fields{"a %= 4", 1},
			args:   args{left: &ast.Atom{Value: `a`, Type: lx.Identifier}, my_op: nil},
			want: &ast.Assignment{
				Identifier: `a`,
				Value: &ast.Operation{
					Type:  lx.ModEq,
					Value: &ast.Atom{Value: `4`, Type: lx.Int},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.fields.input)
			p.curr = tt.fields.curr
			if tt.shouldPanic {
				defer expectPanic(t, nil)
			}
			if got := p.parse_binary(tt.args.left, tt.args.my_op); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parse_binary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_attempt_parse_call(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        core.Expression
		shouldPanic bool
	}{
		{
			name:  `call with two args`,
			input: `print(1, "hello")`,
			want: &ast.Call{
				Name: `print`,
				Args: []core.Expression{&ast.Atom{Type: `int`, Value: `1`}, &ast.Atom{Type: `string`, Value: `hello`}},
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
			name:  `sign-specified integer (minus)`,
			input: `-1`,
			want:  &ast.Atom{Value: `1`, Type: lx.Int, Signed: true},
		},
		{
			name:  `sign-specified exponent float (plus)`,
			input: `+1.2E-7`,
			want:  &ast.Atom{Value: `1.2E-7`, Type: lx.Float},
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
		{
			name:  `negated identifier`,
			input: `!foo`,
			want:  &ast.Atom{Value: `foo`, Type: lx.Identifier, Negated: true},
		},
		{
			name:        `negted parenthesis with non-booleanable expression`,
			input:       `!(v = 5)`,
			shouldPanic: true,
		},
		{
			name:  `negated parenthesis with booleanable expression`,
			input: `!(+x >= 2.8)`,
			want: &ast.Binary{
				Left:     &ast.Atom{Value: `x`, Type: lx.Identifier},
				Operator: lx.Token{Value: `>=`, Type: lx.GreaterThanOrEqual, Line: 1, Column: 6},
				Right:    &ast.Atom{Value: `2.8`, Type: lx.Float},
				Negated:  true,
			},
		},
		{
			name:        `negated parenthesis with non-booleanable expression`,
			input:       `!(x + 2)`,
			shouldPanic: true,
		},
		{
			name:  `negated call`,
			input: `!isWild()`,
			want:  &ast.Call{Name: `isWild`, Negated: true, Args: []core.Expression{}},
		},
		{
			name:  `negated object call`,
			input: `!foo.print(1, "hello")`,
			want: &ast.Call{
				Name:    `print`,
				Args:    []core.Expression{&ast.Atom{Type: `int`, Value: `1`}, &ast.Atom{Type: `string`, Value: `hello`}},
				Object:  `foo`,
				Negated: true,
			},
		},
		{
			name:        `negated literals are not allowed - string`,
			input:       `!"nope"`,
			shouldPanic: true,
		},
		{
			name:        `negated literals are not allowed - float`,
			input:       `!1.234`,
			shouldPanic: true,
		},
		{
			name:        `cannot sign and negate`,
			input:       `-!foo`,
			shouldPanic: true,
		},
		{
			name:        `cannot negate and sign`,
			input:       `!-foo`,
			shouldPanic: true,
		},
		{
			name:  `sign specified identifier`,
			input: `+foo`,
			want:  &ast.Atom{Value: `foo`, Type: lx.Identifier},
		},
		{
			name:  `signed identifier`,
			input: `-foo`,
			want:  &ast.Atom{Value: `foo`, Type: lx.Identifier, Signed: true},
		},
		{
			name:        `signed parenthesis with non-sign-specifiable expression`,
			input:       `-(x >= 2.8)`,
			shouldPanic: true,
		},
		{
			name:        `signed parenthesis with non-sign-specifiable expression`,
			input:       `-(v = 5)`,
			shouldPanic: true,
		},
		{
			name:  `signed parenthesis with sign-specifiable expression`,
			input: `-(x + -y)`,
			want: &ast.Binary{
				Left:     &ast.Atom{Value: `x`, Type: lx.Identifier},
				Operator: lx.Token{Value: `+`, Type: lx.Plus, Line: 1, Column: 5},
				Right:    &ast.Atom{Value: `y`, Type: lx.Identifier, Signed: true},
				Signed:   true,
			},
		},
		{
			name:  `sign specified  call`,
			input: `+isWild()`,
			want:  &ast.Call{Name: `isWild`, Args: []core.Expression{}},
		},
		{
			name:  `signed call`,
			input: `-isWild()`,
			want:  &ast.Call{Name: `isWild`, Signed: true, Args: []core.Expression{}},
		},
		{
			name:  `signed object call`,
			input: `-foo.print(1, "hello")`,
			want: &ast.Call{
				Name:   `print`,
				Args:   []core.Expression{&ast.Atom{Type: `int`, Value: `1`}, &ast.Atom{Type: `string`, Value: `hello`}},
				Object: `foo`,
				Signed: true,
			},
		},
		{
			name:        `signed literals are not allowed - string`,
			input:       `-"nope"`,
			shouldPanic: true,
		},
		{
			name:        `signed literals are not allowed - bool`,
			input:       `-true`,
			shouldPanic: true,
		},
		{
			name:        `signed literals are not allowed - runes`,
			input:       `-'L'`,
			shouldPanic: true,
		},
		{
			name:        `signed literals are not allowed - underscore`,
			input:       `-_`,
			shouldPanic: true,
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
