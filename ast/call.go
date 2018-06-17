package ast

import "github.com/amupitan/hero/ast/core"

type Call struct {
	core.Expression
	Name   string // TODO: take in complete token?
	Args   []core.Expression
	Object string
}

func (c *Call) String() string {
	s := c.Name + `(` + core.StringifyExpressions(c.Args) + `)`
	if c.Object != `` {
		s = c.Object + `.` + s
	}
	return s
}
