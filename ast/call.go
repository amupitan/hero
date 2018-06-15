package ast

import "github.com/amupitan/hero/ast/core"

type Call struct {
	core.Expression
	Name   string // TODO: take in complete token?
	Args   []core.Expression
	object interface{}
}

func (c *Call) String() string {
	return ``
}
