package ast

import "github.com/amupitan/hero/ast/core"

type Block struct {
	core.Expression
	expressions []core.Expression
}