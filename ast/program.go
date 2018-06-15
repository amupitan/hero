package ast

import "github.com/amupitan/hero/ast/core"

type Program struct {
	core.Statement
	Statements []core.Statement
}
