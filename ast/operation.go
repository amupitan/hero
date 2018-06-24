package ast

import (
	"github.com/amupitan/hero/ast/core"
	"github.com/amupitan/hero/lexer"
)

// Operation represents non-binar operations like increment, decrement
// op-equals and possibly delete
type Operation struct {
	core.Expression
	Type lexer.TokenType

	// optional value of the operator if it uses one
	// like +=, ...
	Value core.Expression
}

func (o *Operation) String() string {
	return string(o.Type)
}
