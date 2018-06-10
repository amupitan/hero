package ast

import "github.com/amupitan/hero/ast/core"

type Call struct {
	core.Expression
	name   string
	args   []Param
	object interface{}
}
