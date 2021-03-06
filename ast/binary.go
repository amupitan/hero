package ast

import (
	"github.com/amupitan/hero/ast/core"
	"github.com/amupitan/hero/lexer"
)

type Binary struct {
	core.Expression
	Left     core.Expression
	Right    core.Expression
	Operator lexer.Token
	Negated  bool
	Signed  bool
}

func (b *Binary) String() string {
	return `(` + b.Left.String() + b.Operator.Value + b.Right.String() + `)`
}
