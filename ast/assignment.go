package ast

import (
	"github.com/amupitan/hero/ast/core"
)

type Assignment struct {
	core.Expression
	Identifier string
	Value      core.Expression
}

func (a *Assignment) String() string {
	return a.Identifier + `=` + a.Value.String()
}
