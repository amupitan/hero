package ast

import (
	"testing"

	"github.com/amupitan/hero/types"
)

func TestParam_String(t *testing.T) {
	p := &Param{Name: `firstname`, Type: types.String}

	expects := `firstname string`
	if got := p.String(); got != expects {
		t.Errorf("Param.String() = %s, Expected: %s", got, expects)
	}
}

func TestFunction_String(t *testing.T) {
	f := &Function{
		Definition: Definition{
			Name: `print`,
		},
		Parameters:  []*Param{&Param{Name: `firstname`, Type: types.String}, &Param{Name: `id`, Type: types.Int}},
		Lambda:      false,
		ReturnTypes: []types.Type{types.String, types.Bool},
	}

	expects := `func print(firstname string, id int) (string, bool) {}`
	if got := f.String(); got != expects {
		t.Errorf("Function.String() = %s, Expected: %s", got, expects)
	}
}
