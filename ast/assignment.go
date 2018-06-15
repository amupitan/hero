package ast

import (
	"github.com/amupitan/hero/ast/core"
	"github.com/amupitan/hero/lexer"
)

type Assignment struct {
	Binary
	Identifier lexer.Token
	operator   lexer.Token
	Value      core.Expression
}
