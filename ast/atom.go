package ast

import (
	"github.com/amupitan/hero/ast/core"
	"github.com/amupitan/hero/lexer"
)

type Atom struct {
	lexer.Token
	Negated bool
	Signed bool
}

func (a *Atom) String() string {
	if a.Token.Type == lexer.RawString {
		return `r'` + a.Value
	}
	return a.Value
}

func (a *Atom) Type() core.ExpressionType {
	return core.AtomNode
}
