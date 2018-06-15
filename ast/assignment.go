package ast

import (
	"github.com/amupitan/hero/ast/core"
	"github.com/amupitan/hero/lexer"
)

type Assignment struct {
	core.Expression
	Identifier lexer.Token
	operator   lexer.Token
	Value      core.Expression
}

func (a *Assignment) String() string {
	return ``
}
