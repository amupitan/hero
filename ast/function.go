package ast

import (
	"strings"

	"github.com/amupitan/hero/ast/core"
	"github.com/amupitan/hero/types"
)

type funcBody struct{}

type Param struct {
	Name string
	Type types.Type
}

type Function struct {
	core.Expression
	Definition
	Parameters  []*Param // TODO(DEV) use param,
	Lambda      bool
	ReturnTypes []types.Type
	Body        []core.Statement
	Owner       types.Type
	Private     bool
}

func (p Param) String() string {
	return p.Name + ` ` + p.Type.String()
}

func (f *Function) String() string {
	s := `func ` + f.Name + `(` + stringifyParams(f.Parameters) + `)`
	if len(f.ReturnTypes) > 0 {
		s += ` (` + stringifyTypes(f.ReturnTypes) + `)`
	}
	return s
}

// stringify converts a slice of [Param]s to a comma delimeted string
func stringifyParams(params []*Param) string {
	s := strings.Builder{}
	for i, param := range params {
		s.WriteString(param.String())

		// write comma if not last exp
		if i+1 < len(params) {
			s.WriteString(`, `)
		}
	}

	return s.String()
}

func stringifyTypes(types_ []types.Type) string {
	s := strings.Builder{}
	for i, t := range types_ {
		s.WriteString(t.String())

		// write comma if not last exp
		if i+1 < len(types_) {
			s.WriteString(`, `)
		}
	}

	return s.String()
}
