package ast

import (
	"github.com/amupitan/hero/ast/core"
	"github.com/amupitan/hero/lexer"
)

type Atom struct {
	core.Expression
	Type    lexer.TokenType
	Value   string
	Negated bool
}

func (a *Atom) String() string {
	if a.Type == lexer.RawString {
		return `r'` + a.Value
	}
	return a.Value
}
