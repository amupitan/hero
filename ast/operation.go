package ast

import (
	"github.com/amupitan/hero/ast/core"
	"github.com/amupitan/hero/lexer"
)

// Operation represents non-binar operations like increment, decrement
// and possibly delete
type Operation struct {
	core.Expression
	Type lexer.TokenType
}

func (o *Operation) String() string {
	return string(o.Type)
}
