package ast

import (
	"strings"

	"github.com/amupitan/hero/ast/core"
	lx "github.com/amupitan/hero/lexer"
	"github.com/amupitan/hero/types"
)

type funcBody struct{}

type Param struct {
	Name lx.Token
	Type lx.Token
}

type Function struct {
	Definition
	Parameters  []*Param
	Lambda      bool
	ReturnTypes []lx.Token
	Body        *Block
	Owner       lx.Token
	Private     bool
}

func (p Param) String() string {
	return p.Name.Value + ` ` + p.Type.Value
}

func (f *Function) String() string {
	s := `func ` + f.Name.Value + `(` + stringifyParams(f.Parameters) + `)`
	if len(f.ReturnTypes) > 0 {
		s += ` (` + stringifyTokens(f.ReturnTypes) + `)`
	}
	return s + ` {}`
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

func stringifyTokens(tokens []lx.Token) string {
	s := strings.Builder{}
	for i, t := range tokens {
		s.WriteString(t.Value)

		// write comma if not last exp
		if i+1 < len(tokens) {
			s.WriteString(`, `)
		}
	}

	return s.String()
}

func (f *Function) Type() core.ExpressionType {
	return core.FunctionNode
}

type Param_ struct {
	Name string
	Type types.Type
}

type Function_ struct {
	Definition
	Parameters  []*Param_
	Lambda      bool
	ReturnTypes []types.Type
	Body        *Block
	Owner       types.Type
	Private     bool
}
