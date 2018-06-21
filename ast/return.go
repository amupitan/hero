package ast

import "github.com/amupitan/hero/ast/core"

type Return struct {
	core.Statement
	Values []core.Expression
}

func (r *Return) String() string {
	return `return` + core.StringifyExpressions(r.Values)
}
