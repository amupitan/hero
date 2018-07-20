package ast

import (
	"testing"

	"github.com/amupitan/hero/ast/core"
	"github.com/amupitan/hero/lexer"
	"github.com/amupitan/hero/types"
)

func TestCall_String(t *testing.T) {
	c := &Call{
		Args: []core.Expression{
			&Atom{Token: lexer.Token{Type: lexer.Identifier, Value: `foo`, Line: 1, Column: 1}},
			&Atom{Token: lexer.Token{Type: lexer.Bool, Value: `true`}},
		},
		Name:   `print`,
		Object: `obj`,
	}

	expects := `obj.print(foo, true)`
	if got := c.String(); got != expects {
		t.Errorf("Param.String() = %s, Expected: %s", got, expects)
	}

	// test lambda call
	c = &Call{
		Args: []core.Expression{
			&Atom{Token: lexer.Token{Type: lexer.Identifier, Value: `foo`}},
			&Atom{Token: lexer.Token{Type: lexer.Bool, Value: `true`}},
		},
		Func: &Function{Lambda: true, Parameters: []*Param{}, ReturnTypes: []types.Type{}},
	}

	expects = `func () {}(foo, true)`
	if got := c.String(); got != expects {
		t.Errorf("Param.String() = %s, Expected: %s", got, expects)
	}
}
