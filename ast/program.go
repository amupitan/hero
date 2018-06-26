package ast

import "github.com/amupitan/hero/ast/core"

type Program struct {
	core.Statement
	Body *Block
}

func (p *Program) String() string {
	return `program: ` + p.Body.String()
}
