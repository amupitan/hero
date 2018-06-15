package ast

import "github.com/amupitan/hero/ast/core"

type Block struct {
	core.Statement
	expressions []core.Expression
}

func (b *Block) String() string {
	return `{ ` + core.StringifyExpressions(b.expressions) + `}`
}
