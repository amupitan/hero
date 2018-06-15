package ast

import (
	"github.com/amupitan/hero/ast/core"
	"github.com/amupitan/hero/lexer"
)

type Atom struct {
	core.Expression
	Type  lexer.TokenType
	Value string
}
