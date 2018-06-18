package ast

import "github.com/amupitan/hero/ast/core"

type Block struct {
	core.Statement
	Statements []core.Statement
}

func (b *Block) String() string {
	return `{ ` + core.StringifyStatements(b.Statements) + `}`
}
