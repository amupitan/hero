package ast

import (
	"github.com/amupitan/hero/ast/core"
	"github.com/amupitan/hero/lexer"
)

type Return struct {
	core.Statement
	lexer.Token
	Values map[*lexer.Token]core.Expression
}

func (r *Return) String() string {
	return `return` //+ core.StringifyExpressions(r.Values)
}
