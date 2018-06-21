package ast

import (
	"github.com/amupitan/hero/ast/core"
)

type If struct {
	core.Statement
	// [Condition] represents the condition for evaluation
	// if [Condition] is nil, then it is an else-only
	Condition core.Expression
	// Else represents an else clause
	Else *If
	Body *Block
	// an assignment or definition in an if-block
	Definition core.Expression //TODO(DEV) use this?
}

func (i *If) String() string {
	return `if ` + i.Condition.String() + `{}`
}
