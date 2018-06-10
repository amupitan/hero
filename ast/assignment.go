package ast

import (
	"github.com/amupitan/hero/ast/core"
	"github.com/amupitan/hero/lexer"
)

type Assignment struct {
	core.Expression
	identifier lexer.Token
	operator   lexer.Token
	value      lexer.Token
}
